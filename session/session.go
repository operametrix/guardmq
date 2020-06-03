package session

import (
	"net"
)

type Session struct {
	InboundConn   net.Conn
	OutboundConn  net.Conn
	ClientID      string
	Username      string
	Password      []byte
	KeepAlive     uint16
	IsDestroy     bool
	Subscriptions []string
}

func (session *Session) Destroy() {
	if session.IsDestroy == false {
		session.IsDestroy = true
	}
}
