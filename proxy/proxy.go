package proxy

import (
	"github.com/eclipse/paho.mqtt.golang/packets"
	"operametrix/mqtt/session"
	"net"
)

func ForwardToBroker(current_session *session.Session, packet *packets.ControlPacket) {
	(*packet).Write(current_session.OutboundConn)
}

func ForwardToClient(current_session *session.Session, packet *packets.ControlPacket) {
	(*packet).Write(current_session.InboundConn)
}

func SocketReader(conn net.Conn, packetChannel chan packets.ControlPacket, errorChannel chan error) {
	for {
		packet, err := packets.ReadPacket(conn)
		if err != nil {
			errorChannel <- err
			return
		}

		packetChannel <- packet
	}
}