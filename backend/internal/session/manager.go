package session

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

const (
	SessionCookie    = "SESSION"
	SessionObjectKey = "user"
)

type Manager interface {
	CreateNewSession(session UserData, r *http.Request, w http.ResponseWriter) error
	CheckSession(r *http.Request, w http.ResponseWriter) (session UserData, ok bool, err error)
}

type scsManager struct {
	scsManager *scs.SessionManager
}

func NewManager(manager *scs.SessionManager) Manager {
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
