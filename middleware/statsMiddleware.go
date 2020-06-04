package middleware

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"operametrix/mqtt/session"
	"log"
	"time"
)

func StatsRoutineMiddleware() {
	for {
		time.Sleep(5 * time.Second)
		log.Println("*** STATS")
		for _, s := range session.ActiveSessions {
			s.Stats.Mux.Lock()
			log.Println(s.InboundConn.RemoteAddr(), "IN:", s.Stats.InTotalPackets*12, "pkt/min OUT:", s.Stats.OutTotalPackets*12, "pkt/min")
			s.Stats.InTotalPackets  = 0
			s.Stats.OutTotalPackets = 0
			s.Stats.Mux.Unlock()
		}
		log.Println("***")
	}
}

func StatsOutMiddleware(next Handler) Handler {
	return HandlerFunc(func(current_session *session.Session, packet *packets.ControlPacket) {
		current_session.Stats.Mux.Lock()
		current_session.Stats.OutTotalPackets += 1

		switch (*packet).(type) {
		case *packets.PublishPacket:
			current_session.Stats.OutPublishPackets += 1
		}

		current_session.Stats.Mux.Unlock()
		next.Serve(current_session, packet)
	})
}

func StatsInMiddleware(next Handler) Handler {
	return HandlerFunc(func(current_session *session.Session, packet *packets.ControlPacket) {
		current_session.Stats.Mux.Lock()
		current_session.Stats.InTotalPackets += 1

		switch (*packet).(type) {
		case *packets.PublishPacket:
			current_session.Stats.InPublishPackets += 1
		}
		current_session.Stats.Mux.Unlock()
		next.Serve(current_session, packet)
	})
}