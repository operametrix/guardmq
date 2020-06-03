package middleware

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"log"
	"operametrix/mqtt/session"
)

func ActiveMiddleware(next Handler) Handler {
	return HandlerFunc(func(current_session *session.Session, packet *packets.ControlPacket) {

		switch p := (*packet).(type) {
		case *packets.ConnectPacket:
			current_session.ClientID = p.ClientIdentifier
			current_session.KeepAlive = p.Keepalive
			if p.UsernameFlag {
				current_session.Username = p.Username
			}
			if p.PasswordFlag {
				current_session.Password = p.Password
			}

			connack := packets.NewControlPacket(packets.Connack)
			connack.Write(current_session.OutboundConn)

		case *packets.ConnackPacket:

		case *packets.DisconnectPacket:
			current_session.Destroy()

		case *packets.SubscribePacket:
			for _, topic := range p.Topics {
				log.Println("Add subscription:", topic)
				current_session.Subscriptions = append(current_session.Subscriptions, topic)
			}

			suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
			suback.MessageID = p.MessageID
			suback.ReturnCodes = append(suback.ReturnCodes, 0)
			suback.Write(current_session.OutboundConn)

		case *packets.SubackPacket:

		case *packets.UnsubscribePacket:

		case *packets.UnsubackPacket:

		case *packets.PublishPacket:
			log.Println("Publish:", p.Payload)

		case *packets.PubackPacket:

		case *packets.PubrecPacket:

		case *packets.PubrelPacket:

		case *packets.PubcompPacket:

		case *packets.PingreqPacket:
			pingresp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
			pingresp.Write(current_session.OutboundConn)

		case *packets.PingrespPacket:

		}

		next.Serve(current_session, packet)
	})
}
