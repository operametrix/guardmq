package session

import (
	"log"
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

func (session *Session) Print() {
	log.Println("New session from clientID:", session.ClientID)
}

func (session *Session) Destroy() {
	if session.IsDestroy == false {
		log.Println("Destroy session from clientID:", session.ClientID)
		session.IsDestroy = true
	}
}
