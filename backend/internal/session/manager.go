package session

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

type UserData struct {
	UserId uint64
}

func init() {
	gob.Register(UserData{})
}

type Manager interface {
	CreateNewSession(session UserData, r *http.Request, w http.ResponseWriter) error
	CheckSession(r *http.Request, w http.ResponseWriter) (session UserData, ok bool, err error)
}

type scsManager struct {
	scsManager *scs.SessionManager
}

func newManager(manager *scs.SessionManager) Manager {
	return &scsManager{manager}
}

func (m *scsManager) CreateNewSession(session UserData, r *http.Request, w http.ResponseWriter) error {
	if err := m.scsManager.RenewToken(r.Context()); err != nil {
		return fmt.Errorf("creating new session: %w", err)
	}
	m.scsManager.Put(r.Context(), SessionObjectKey, session)
	return nil
}

func (m *scsManager) CheckSession(r *http.Request, w http.ResponseWriter) (session UserData, ok bool, err error) {
	session, ok = m.scsManager.Get(r.Context(), SessionObjectKey).(UserData)
	return
}

func newScsRedisStore(redisClient *redis.Client) scs.Store {
	return goredisstore.New(redisClient)
}

func newScsManager(store scs.Store) *scs.SessionManager {
	scsManager := scs.New()
	scsManager.Store = store
	scsManager.Cookie.Name = SessionCookie
	scsManager.Cookie.HttpOnly = true
	scsManager.Cookie.Secure = true
	scsManager.Cookie.Persist = true
	scsManager.IdleTimeout = 60 * time.Hour
	return scsManager
}
