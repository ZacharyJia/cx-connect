package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type AdminSessionGroup struct {
	SessionKey      string                `json:"session_key"`
	Platform        string                `json:"platform"`
	ActiveSessionID string                `json:"active_session_id"`
	Interactive     bool                  `json:"interactive"`
	Sessions        []AdminSessionSummary `json:"sessions"`
}

type AdminSessionSummary struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	WorkDir        string    `json:"work_dir"`
	AgentSessionID string    `json:"agent_session_id"`
	HistoryCount   int       `json:"history_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Busy           bool      `json:"busy"`
	Active         bool      `json:"active"`
}

type AdminSessionDetail struct {
	Project         string          `json:"project"`
	SessionKey      string          `json:"session_key"`
	Platform        string          `json:"platform"`
	ActiveSessionID string          `json:"active_session_id"`
	Interactive     bool            `json:"interactive"`
	Session         SessionSnapshot `json:"session"`
}

type AdminCreateSessionResult struct {
	Project         string          `json:"project"`
	SessionKey      string          `json:"session_key"`
	Platform        string          `json:"platform"`
	ActiveSessionID string          `json:"active_session_id"`
	Interactive     bool            `json:"interactive"`
	Session         SessionSnapshot `json:"session"`
	DisplayWorkDir  string          `json:"display_work_dir"`
}

func (e *Engine) AdminSessionGroups() []AdminSessionGroup {
	interactive := make(map[string]bool)

	e.interactiveMu.Lock()
	for key, state := range e.interactiveStates {
		if state != nil && state.platform != nil {
			interactive[key] = true
		}
	}
	e.interactiveMu.Unlock()

	users := e.sessions.Snapshot()
	groups := make([]AdminSessionGroup, 0, len(users))
	for _, user := range users {
		summaries := make([]AdminSessionSummary, 0, len(user.Sessions))
		for _, session := range user.Sessions {
			summaries = append(summaries, AdminSessionSummary{
				ID:             session.ID,
				Name:           session.Name,
				WorkDir:        session.WorkDir,
				AgentSessionID: session.AgentSessionID,
				HistoryCount:   session.HistoryCount,
				CreatedAt:      session.CreatedAt,
				UpdatedAt:      session.UpdatedAt,
				Busy:           session.Busy,
				Active:         session.ID == user.ActiveSessionID,
			})
		}
		sort.Slice(summaries, func(i, j int) bool {
			return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
		})
		groups = append(groups, AdminSessionGroup{
			SessionKey:      user.UserKey,
			Platform:        platformNameFromSessionKey(user.UserKey),
			ActiveSessionID: user.ActiveSessionID,
			Interactive:     interactive[user.UserKey],
			Sessions:        summaries,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		left := time.Time{}
		right := time.Time{}
		if len(groups[i].Sessions) > 0 {
			left = groups[i].Sessions[0].UpdatedAt
		}
		if len(groups[j].Sessions) > 0 {
			right = groups[j].Sessions[0].UpdatedAt
		}
		if left.Equal(right) {
			return groups[i].SessionKey < groups[j].SessionKey
		}
		return left.After(right)
	})

	return groups
}

func (e *Engine) AdminSessionDetail(sessionKey, sessionID string) (AdminSessionDetail, error) {
	session, ok := e.sessions.SessionSnapshot(sessionKey, sessionID, true)
	if !ok {
		return AdminSessionDetail{}, fmt.Errorf("session %q not found", sessionID)
	}

	e.interactiveMu.Lock()
	state := e.interactiveStates[sessionKey]
	interactive := state != nil && state.platform != nil
	e.interactiveMu.Unlock()

	return AdminSessionDetail{
		Project:         e.name,
		SessionKey:      sessionKey,
		Platform:        platformNameFromSessionKey(sessionKey),
		ActiveSessionID: e.sessions.ActiveSessionID(sessionKey),
		Interactive:     interactive,
		Session:         session,
	}, nil
}

func (e *Engine) SubmitPrompt(sessionKey, sessionID, prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("prompt is required")
	}

	session, ok := e.sessions.GetSession(sessionKey, sessionID)
	if !ok {
		return fmt.Errorf("session %q not found", sessionID)
	}

	platform, replyCtx, err := e.resolvePromptContext(sessionKey)
	if err != nil {
		return err
	}

	currentActiveID := e.sessions.ActiveSessionID(sessionKey)
	if currentActiveID != session.ID {
		if _, err := e.sessions.SetActiveSession(sessionKey, session.ID); err != nil {
			return err
		}
		e.cleanupInteractiveState(sessionKey)
	}

	if !session.TryLock() {
		return fmt.Errorf("session %q is busy", session.Name)
	}

	msg := &Message{
		SessionKey: sessionKey,
		Platform:   platform.Name(),
		UserID:     "web",
		UserName:   "web",
		Content:    prompt,
		ReplyCtx:   replyCtx,
	}
	go e.processInteractiveMessage(platform, msg, session)
	return nil
}

func (e *Engine) CreateSession(sessionKey, name, workDirInput string) (AdminCreateSessionResult, error) {
	sessionKey = strings.TrimSpace(sessionKey)
	if sessionKey == "" {
		sessionKey = newWebSessionKey()
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = newRandomSessionName()
	}

	if strings.TrimSpace(workDirInput) == "" {
		workDirInput = "default"
	}

	workDir, displayDir, err := e.resolveSessionWorkDir(name, workDirInput)
	if err != nil {
		return AdminCreateSessionResult{}, err
	}

	session := e.sessions.NewSession(sessionKey, name, workDir)
	return AdminCreateSessionResult{
		Project:         e.name,
		SessionKey:      sessionKey,
		Platform:        platformNameFromSessionKey(sessionKey),
		ActiveSessionID: session.ID,
		Interactive:     false,
		Session:         session.Snapshot(true),
		DisplayWorkDir:  displayDir,
	}, nil
}

func (e *Engine) resolvePromptContext(sessionKey string) (Platform, any, error) {
	e.interactiveMu.Lock()
	state := e.interactiveStates[sessionKey]
	e.interactiveMu.Unlock()

	if state != nil {
		state.mu.Lock()
		platform := state.platform
		replyCtx := state.replyCtx
		state.mu.Unlock()
		if platform != nil && replyCtx != nil {
			return platform, replyCtx, nil
		}
	}

	platformName := platformNameFromSessionKey(sessionKey)
	if platformName == "" {
		return nil, nil, fmt.Errorf("invalid session key %q", sessionKey)
	}

	platform := e.findPlatform(platformName)
	if platform == nil {
		return nil, nil, fmt.Errorf("platform %q not found", platformName)
	}

	reconstructor, ok := platform.(ReplyContextReconstructor)
	if !ok {
		return nil, nil, fmt.Errorf("platform %q does not support prompt injection without an active session", platformName)
	}

	replyCtx, err := reconstructor.ReconstructReplyCtx(sessionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("reconstruct reply context: %w", err)
	}
	return platform, replyCtx, nil
}

func (e *Engine) findPlatform(name string) Platform {
	if name == "web" && e.webPlatform != nil {
		return e.webPlatform
	}
	for _, platform := range e.platforms {
		if platform.Name() == name {
			return platform
		}
	}
	return nil
}

func platformNameFromSessionKey(sessionKey string) string {
	head, _, found := strings.Cut(sessionKey, ":")
	if !found {
		return ""
	}
	return head
}

func newWebSessionKey() string {
	var b [6]byte
	if _, err := rand.Read(b[:]); err == nil {
		return "web:" + hex.EncodeToString(b[:])
	}
	return fmt.Sprintf("web:%d", time.Now().UnixNano())
}
