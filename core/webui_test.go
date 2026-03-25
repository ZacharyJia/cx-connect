package core

import (
	"strings"
	"testing"
)

func TestAdminWebUIKeepsSessionListsScrollable(t *testing.T) {
	requiredSnippets := []string{
		".session-groups {",
		"overflow-y: auto;",
		".group-sessions {",
		"max-height: min(52vh, 420px);",
		"'<div class=\"group-sessions\">' + sessions + '</div>' +",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(adminWebUIHTML, snippet) {
			t.Fatalf("admin web ui is missing scroll guard %q", snippet)
		}
	}
}
