package cli

import (
	"context"

	"github.com/nikoksr/proji/internal/config"
	"github.com/nikoksr/proji/pkg/packages"
	"github.com/nikoksr/proji/pkg/projects"
)

// Session is meant to provide commonly used functionality for CLI sessions. Another suitable name for this type would
// be `Context`, but we went with `Session` to avoid confusion with `context.Context`.
// Session is primarily meant to be used as a context.Context value and passed to commands. At the time of writing,
// there is no pretty way to do this. And since we don't want to add the logic for creating core types to every command,
// we decided to create them in the root command and pass them to the subcommands via context. Yes, this is probably
// (hopefully?) not the prettiest way to do it, but it works.
// Due to the fact that we create the session in the root command, only absolutely essential types are embedded in the
// session. A manager.PackageManager is not necessary for every command, thus we don't embed it. The loaded app config
// on the other hand is necessary for almost every command, so we embed it. The session type may be extended in the
// future to include more types.
// A logging.Logger instance is not embedded since we use logging frequently outside a CLI's session, thus, to avoid
// confusion, we don't embed it and instead handle loggers manually.
type Session struct {
	Debug          bool
	Config         *config.Config
	PackageManager packages.Manager
	ProjectManager projects.Manager
}

// NewSessionWithMode creates a new session with the given debug mode.
func NewSessionWithMode(debug bool) *Session {
	return &Session{Debug: debug}
}

// In const variable to make it easier for testing purposes.
const defaultDebug = false

// NewSession creates a new session. It uses default values. By default, the session is not in debug mode.
func NewSession() *Session {
	return NewSessionWithMode(defaultDebug)
}

// WithConfig sets the given config on the session. It returns the session to allow chaining.
func (session *Session) WithConfig(config *config.Config) *Session {
	session.Config = config

	return session
}

// WithPackageManager sets the given package manager on the session. It returns the session to allow chaining.
func (session *Session) WithPackageManager(manager packages.Manager) *Session {
	session.PackageManager = manager

	return session
}

// WithProjectManager sets the given project manager on the session. It returns the session to allow chaining.
func (session *Session) WithProjectManager(manager projects.Manager) *Session {
	session.ProjectManager = manager

	return session
}

// As recommended by 'revive' linter.
type contextKey string

const sessionKey contextKey = "session"

// WithSession returns a new context.Context with the given session. The session is not copied, so it should not be
// modified after it is passed to this function.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// SessionFromContext returns the session from the given context. If the context does not contain a session, the default
// session is returned.
func SessionFromContext(ctx context.Context) *Session {
	if ctx == nil {
		return NewSession()
	}

	if session, ok := ctx.Value(sessionKey).(*Session); ok {
		return session
	}

	return NewSession()
}
