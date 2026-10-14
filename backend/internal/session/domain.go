package session

import "encoding/gob"

type UserData struct {
	UserId uint64
}

func init() {
	gob.Register(UserData{})
}
