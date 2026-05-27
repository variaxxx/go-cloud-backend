package auth_jwt

import (
	"strings"
	"testing"
	"time"
)

func TestManagerIssueAndParse(t *testing.T) {
	manager := NewManager(config{Secret: "test-secret", Ttl: time.Hour})

	token, err := manager.Issue(42)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if token == "" {
		t.Fatal("Issue() token is empty")
	}

	userID, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("Parse() userID = %d, want 42", userID)
	}
}

func TestManagerParseRejectsInvalidToken(t *testing.T) {
	manager := NewManager(config{Secret: "test-secret", Ttl: time.Hour})

	_, err := manager.Parse("not-a-jwt")
	if err == nil {
		t.Fatal("Parse() error = nil, want parse error")
	}
	if !strings.Contains(err.Error(), "parse token") {
		t.Fatalf("Parse() error = %v, want parse token context", err)
	}
}
