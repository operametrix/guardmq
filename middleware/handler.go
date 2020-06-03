package middleware

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"operametrix/mqtt/session"
)

type Handler interface {
	Serve(current_session *session.Session, packet *packets.ControlPacket)
}

type HandlerFunc func(*session.Session, *packets.ControlPacket)

func (f HandlerFunc) Serve(current_session *session.Session, packet *packets.ControlPacket) {
	f(current_session, packet)
}
