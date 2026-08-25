package security

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

const (
	SessionCookie    = "SESSION"
	SessionObjectKey = "user"
)

type UserSession struct {
	UserId uint64
}

func init() {
	gob.Register(UserSession{})
}

type SessionManager interface {
	CreateNewSession(session UserSession, r *http.Request, w http.ResponseWriter) error
	CheckSession(r *http.Request, w http.ResponseWriter) (session UserSession, ok bool, err error)
}

type ScsSessionManager struct {
	scsManager *scs.SessionManager
}

func NewAppSessionManager(manager *scs.SessionManager) *ScsSessionManager {
	return &ScsSessionManager{manager}
}

func (m *ScsSessionManager) CreateNewSession(session UserSession, r *http.Request, w http.ResponseWriter) error {
	if err := m.scsManager.RenewToken(r.Context()); err != nil {
		return fmt.Errorf("creating new session: %w", err)
	}
	m.scsManager.Put(r.Context(), SessionObjectKey, session)
	return nil
}

func (m *ScsSessionManager) CheckSession(r *http.Request, w http.ResponseWriter) (session UserSession, ok bool, err error) {
	session, ok = m.scsManager.Get(r.Context(), SessionObjectKey).(UserSession)
	return
}

func NewScsRedisSessionStore(redisClient *redis.Client) scs.Store {
	return goredisstore.New(redisClient)
}

func NewScsSessionManager(store scs.Store) *scs.SessionManager {
	scsManager := scs.New()
	scsManager.Store = store
	scsManager.Cookie.Name = SessionCookie
	scsManager.Cookie.HttpOnly = true
    scsManager.Cookie.Secure = true
	scsManager.Cookie.Persist = true
	scsManager.IdleTimeout = 60 * time.Hour
	return scsManager
}
