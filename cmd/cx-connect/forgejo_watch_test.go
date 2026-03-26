package main

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZacharyJia/cx-connect/config"
	"github.com/ZacharyJia/cx-connect/core"
	"github.com/ZacharyJia/cx-connect/forgejowatch"
)

func TestEngineAdminClientCreateSessionAndListGroups(t *testing.T) {
	workDir := t.TempDir()
	engine := core.NewEngine(defaultProjectName, &stubAgent{}, nil, filepath.Join(t.TempDir(), "sessions.json"), core.LangEnglish, []config.AllowUser{})
	client := newEngineAdminClient(engine.ProjectName(), engine)

	created, err := client.CreateSession(context.Background(), forgejowatch.CreateSessionRequest{
		SessionKey: "web:test-user",
		Name:       "forgejo-ops-12",
		WorkDir:    workDir,
	})
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if created.SessionKey != "web:test-user" {
		t.Fatalf("CreateSession() sessionKey = %q", created.SessionKey)
	}
	if created.Session.ID == "" {
		t.Fatal("CreateSession() returned empty session ID")
	}

	groups, err := client.ListSessionGroups(context.Background(), "")
	if err != nil {
		t.Fatalf("ListSessionGroups() error = %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("ListSessionGroups() group count = %d, want 1", len(groups))
	}
	if groups[0].SessionKey != "web:test-user" {
		t.Fatalf("ListSessionGroups() sessionKey = %q", groups[0].SessionKey)
	}
	if len(groups[0].Sessions) != 1 {
		t.Fatalf("ListSessionGroups() session count = %d, want 1", len(groups[0].Sessions))
	}
	if groups[0].Sessions[0].ID != created.Session.ID {
		t.Fatalf("ListSessionGroups() session ID = %q, want %q", groups[0].Sessions[0].ID, created.Session.ID)
	}
}

func TestEngineAdminClientRejectsUnknownProject(t *testing.T) {
	engine := core.NewEngine(defaultProjectName, &stubAgent{}, nil, filepath.Join(t.TempDir(), "sessions.json"), core.LangEnglish, []config.AllowUser{})
	client := newEngineAdminClient(engine.ProjectName(), engine)

	_, err := client.ListSessionGroups(context.Background(), "other-project")
	if err == nil {
		t.Fatal("ListSessionGroups() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), `project "other-project" not found`) {
		t.Fatalf("ListSessionGroups() error = %v", err)
	}
}

type stubAgent struct{}

func (a *stubAgent) Name() string { return "stub-agent" }

func (a *stubAgent) StartSession(ctx context.Context, sessionID string) (core.AgentSession, error) {
	return nil, nil
}

func (a *stubAgent) ListSessions(ctx context.Context) ([]core.AgentSessionInfo, error) {
	return nil, nil
}

func (a *stubAgent) Stop() error { return nil }
