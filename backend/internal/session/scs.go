package session

import (
	"time"

	"github.com/alexedwards/scs/v2"
)

func NewScsManager(store scs.Store) *scs.SessionManager {
	scsManager := scs.New()
	scsManager.Store = store
	scsManager.Cookie.Name = SessionCookie
	scsManager.Cookie.HttpOnly = true
	scsManager.Cookie.Secure = true
	scsManager.Cookie.Persist = true
	scsManager.IdleTimeout = 60 * time.Hour
	return scsManager
}
