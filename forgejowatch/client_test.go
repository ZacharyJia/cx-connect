package forgejowatch

import (
	"encoding/json"
	"testing"
)

func TestForgejoIssueUnmarshalRepositoryOwnerString(t *testing.T) {
	payload := []byte(`{
		"number": 12,
		"title": "Fix login bug",
		"repository": {
			"full_name": "nselab/demo",
			"name": "demo",
			"owner": "nselab"
		}
	}`)

	var issue ForgejoIssue
	if err := json.Unmarshal(payload, &issue); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if issue.Repository.Owner.Login != "nselab" {
		t.Fatalf("unexpected owner: %+v", issue.Repository.Owner)
	}
}

func TestForgejoIssueUnmarshalRepositoryOwnerObject(t *testing.T) {
	payload := []byte(`{
		"number": 13,
		"title": "Fix login bug",
		"repository": {
			"full_name": "nselab/demo",
			"name": "demo",
			"owner": {
				"login": "nselab"
			}
		}
	}`)

	var issue ForgejoIssue
	if err := json.Unmarshal(payload, &issue); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if issue.Repository.Owner.Login != "nselab" {
		t.Fatalf("unexpected owner: %+v", issue.Repository.Owner)
	}
}
