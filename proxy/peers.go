package proxy

import (
    "crypto/tls"
	"crypto/x509"
	"math/rand"
	"io/ioutil"
	"fmt"
	"time"
	"log"
	"net"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

func randomClientID(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
 
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}

type Peer struct {
	Hostname string    `yaml:"hostname"`
	Port     int       `yaml:"port"`
	TLS      bool      `yaml:"tls"`
	CAFile   string    `yaml:"cafile"`
	MTLS     bool      `yaml:"mtls"`
	CertFile string    `yaml:"certfile"`
	KeyFile  string    `yaml:"keyfile"`
	Import   []string  `yaml:"import"`
	Export   []string  `yaml:"export"`
}

func (peer *Peer) Serve() {

	host := fmt.Sprintf("%s:%d", peer.Hostname, peer.Port)
	var inboundConn net.Conn
	var err error

	for {
		err = nil

		if peer.TLS {
			var caCert []byte
			var cert tls.Certificate
			tlsConfig := &tls.Config{}

			if peer.CAFile == "" {
				log.Println(host, ": you must indicate a CA file for TLS")
				return
			}

			caCertPool  := x509.NewCertPool()
			caCert, err = ioutil.ReadFile(peer.CAFile)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool

			if peer.MTLS {
				if peer.CertFile == "" || peer.KeyFile == "" {
					log.Println(host, ": you must indicate a certificate file and a key file for mTLS")
					return
				}

				cert, err = tls.LoadX509KeyPair(peer.CertFile, peer.KeyFile)
				if err != nil {
					log.Fatal(err.Error())
				}
				tlsConfig.Certificates = []tls.Certificate{cert}
			}

			inboundConn, err = tls.Dial("tcp", host, tlsConfig)
			if err != nil {
				log.Println("Failed to contact the peer", host, "with TLS")
				log.Println(err.Error())
			}
		} else {
			inboundConn, err = net.Dial("tcp", host)
			if err != nil {
				log.Println("Failed to contact the peer", host)
			}
		}

		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		outboundConn, err := ConnectLocalBroker()
		if err != nil {
			log.Println("Session with the peer", host, "failed")
			time.Sleep(5 * time.Second)
			continue
		}

		connectPacket := packets.ConnectPacket{
			FixedHeader: packets.FixedHeader{
				MessageType: packets.Connect,
			},
			CleanSession: true,
			ClientIdentifier: randomClientID(20),
			Keepalive: 0,
			ProtocolName: "MQTT",
			ProtocolVersion: 4,
		}
		connectPacket.Write(inboundConn)
		_, err = packets.ReadPacket(inboundConn)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		connectPacket.Write(outboundConn)
		_, err = packets.ReadPacket(outboundConn)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		subPacket := packets.SubscribePacket{
			FixedHeader: packets.FixedHeader{
				MessageType: packets.Subscribe,
				Qos: 1,
			},
			MessageID: 1,
			Topics: peer.Import,
			Qoss: make([]byte,len(peer.Import)),
		}
		subPacket.Write(inboundConn)
		_, err = packets.ReadPacket(inboundConn)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		subPacket = packets.SubscribePacket{
			FixedHeader: packets.FixedHeader{
				MessageType: packets.Subscribe,
				Qos: 1,
			},
			MessageID: 1,
			Topics: peer.Export,
			Qoss: make([]byte,len(peer.Export)),
		}
		subPacket.Write(outboundConn)
		_, err = packets.ReadPacket(outboundConn)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		if peer.TLS {
			log.Println("Open TLS peering session with", host)
		} else {
			log.Println("Open peering session with", host)
		}

		HandleConnection(inboundConn, outboundConn)
		log.Println("Connection failed for peer", host)
		time.Sleep(5 * time.Second)
	}
}
