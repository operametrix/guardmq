package middleware

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"log"
	"operametrix/mqtt/session"
)

func LoggingMiddleware(next Handler) Handler {
	return HandlerFunc(func(current_session *session.Session, packet *packets.ControlPacket) {
		log.Println((*packet).String())
		next.Serve(current_session, packet)
	})
}
