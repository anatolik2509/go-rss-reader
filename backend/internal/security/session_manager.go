package security

import (
	"fmt"
	"net/http"
	"rss-reader-backend/internal/config"

	"github.com/boj/redistore/v2"
	"github.com/gorilla/sessions"
)

const (
	SessionCookie    = "SESSION"
	SessionObjectKey = "user"
)

type UserSession struct {
	userId uint64
}

type SessionManager interface {
	CreateNewSession(session UserSession, r *http.Request, w http.ResponseWriter) error
	CheckSession(r *http.Request, w http.ResponseWriter) (session UserSession, ok bool, err error)
}

type GorillaSessionManager struct {
	store sessions.Store
}

func NewGorillaSessionManager(store sessions.Store) *GorillaSessionManager {
	return &GorillaSessionManager{store: store}
}

func (m *GorillaSessionManager) CreateNewSession(session UserSession, r *http.Request, w http.ResponseWriter) error {
	gorillaSession, err := m.store.Get(r, SessionCookie)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	gorillaSession.Values[SessionObjectKey] = session
	err = m.store.Save(r, w, gorillaSession)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	return nil
}

func (m *GorillaSessionManager) CheckSession(r *http.Request, w http.ResponseWriter) (session UserSession, ok bool, err error) {
	gorillaSession, err := m.store.Get(r, SessionCookie)
	if err != nil {
		return UserSession{}, false, fmt.Errorf("failed to parse db config: %w", err)
	}
	session = gorillaSession.Values[SessionObjectKey].(UserSession)
	ok = !gorillaSession.IsNew
	return
}

func NewGorillaSessionStore(redisConfig config.RedisConfig, secretKey string) (sessions.Store, error) {
	store, err := redistore.NewStore(
		redistore.KeysFromStrings(secretKey),
		redistore.WithAddress("tcp", redisConfig.Host),
		redistore.WithAuth(redisConfig.User, redisConfig.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("creating gorilla session store: %w", err)
	}
	return store, nil
}
