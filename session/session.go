package session

import (
	"net"
	"sync"
)

var (
	ActiveSessions []*Session
)

type SessionStatistics struct {
	Mux sync.Mutex
	InTotalPackets int
	InPublishPackets int
	OutTotalPackets int
	OutPublishPackets int
}

type Session struct {
	CommonName    string
	EndPoint      string
	InboundConn   net.Conn
	OutboundConn  net.Conn
	Stats         SessionStatistics
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
