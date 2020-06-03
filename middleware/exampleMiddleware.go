package middleware

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"log"
	"operametrix/mqtt/session"
)

func ExampleMiddleware(next Handler) Handler {
	return HandlerFunc(func(current_session *session.Session, packet *packets.ControlPacket) {
		log.Println("--- BEFORE PACKET")
		next.Serve(current_session, packet)
		log.Println("--- AFTER PACKET")
	})
}
