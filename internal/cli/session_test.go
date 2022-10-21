package cli

import (
	"context"
	"testing"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/packages"
	"github.com/nikoksr/proji/pkg/projects"
)

func TestNewSession(t *testing.T) {
	t.Parallel()

	session := NewSession()
	if session == nil {
		t.Fatal("NewSession() returned nil")
	}
	if session.Debug != defaultDebug {
		t.Errorf("NewSession() returned a session with debug mode %v, expected %v", session.Debug, defaultDebug)
	}

	session = NewSessionWithMode(true)
	if session == nil {
		t.Fatal("NewSessionWithMode() returned nil")
	}
	if session.Debug != true {
		t.Errorf("NewSessionWithMode() returned a session with debug mode %v, expected %v", session.Debug, true)
	}
}

func TestSessionContext(t *testing.T) {
	t.Parallel()

	session := NewSession()
	if session == nil {
		t.Fatal("NewSession() returned nil")
	}

	ctx := WithSession(context.Background(), session)
	if ctx == nil {
		t.Fatal("WithSession() returned nil")
	}

	// Session
	sessionFromContext := SessionFromContext(ctx)
	if sessionFromContext == nil {
		t.Fatal("SessionFromContext() returned nil")
	}
	if sessionFromContext != session {
		t.Errorf("SessionFromContext() returned %v, expected %v", sessionFromContext, session)
	}

	// Package manager
	var pama packages.Manager
	session.WithPackageManager(pama)
	if session.PackageManager != pama {
		t.Errorf("Session.WithPackageManager() set %v, expected %v", session.PackageManager, pama)
	}

	// Project manager
	var prma projects.Manager
	session.WithProjectManager(prma)
	if session.ProjectManager != prma {
		t.Errorf("Session.WithProjectManager() set %v, expected %v", session.ProjectManager, prma)
	}

	// Nil context
	sessionFromContext = SessionFromContext(nil) //nolint:staticcheck
	if sessionFromContext == nil {
		t.Fatal("SessionFromContext() returned nil")
	}

	// Bad context key
	ctx = context.WithValue(context.Background(), "session", "bad key") //nolint:revive,staticcheck
	sessionFromContext = SessionFromContext(ctx)
	if sessionFromContext == nil {
		t.Errorf("SessionFromContext() returned nil")
	}
}

func TestSession_WithConfig(t *testing.T) {
	t.Parallel()

	session := NewSession()
	if session == nil {
		t.Fatal("NewSession() returned nil")
	}

	conf := &config.Config{}
	session.WithConfig(conf)
	if session.Config != conf {
		t.Errorf("Session.WithConfig() set %v, expected %v", session.Config, conf)
	}
}
