package proxy

import (
	"crypto/x509"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"operametrix/mqtt/session"
)

type Listener struct {
	Hostname string   `yaml:"hostname"`
	Port     int      `yaml:"port"`
	TLS      bool     `yaml:"tls"`
	CertFile string   `yaml:"certfile"`
	KeyFile  string   `yaml:"keyfile"`
	MTLS     bool     `yaml:"mtls"`
	CAFile   string   `yaml:"cafile"`
}

func (listener *Listener) Serve() {
	host := fmt.Sprintf("%s:%s", listener.Hostname, strconv.Itoa(listener.Port))

	var l net.Listener
	var err error
	var commonName string

	if listener.TLS {

		if listener.CertFile == "" || listener.KeyFile == "" {
			log.Println(host, ": you must indicate a certificate file and a key file")
			return
		}

		cert, err := tls.LoadX509KeyPair(listener.CertFile, listener.KeyFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if listener.MTLS {
			if listener.CAFile == "" {
				log.Println(host, ": you must indicate a CA file for mTLS")
				return
			}

			caCertPool  := x509.NewCertPool()
			caCert, err := ioutil.ReadFile(listener.CAFile)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs  = caCertPool
		}

		log.Println("Start to listen on", host, "with TLS")
		l, err = tls.Listen("tcp", host, tlsConfig)
	} else {
		log.Println("Start to listen on", host)
		l, err = net.Listen("tcp", host)
	}

	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("[", host, "] new connection from", conn.RemoteAddr())

		if listener.TLS {
			err := conn.(*tls.Conn).Handshake()
			if err != nil {
				log.Println(err)
				conn.Close()
				continue
			}

			state := conn.(*tls.Conn).ConnectionState()
			if (len(state.VerifiedChains) <= 0) || (len(state.VerifiedChains[0]) <= 0) {
				log.Println("Unverified certificate chain")
				conn.Close()
				continue
			}

			commonName = state.VerifiedChains[0][0].Subject.CommonName
		}

		var current_session session.Session
		current_session.InboundConn = conn
		current_session.EndPoint = conn.RemoteAddr().String()

		if listener.TLS {
			current_session.CommonName = commonName
		} else {
			current_session.CommonName = conn.RemoteAddr().String()
		}

		go LaunchConnection(&current_session)
	}
}

func LaunchConnection(current_session *session.Session) {
	err := ConnectLocalBroker(current_session)
	if err != nil {
		log.Println("Failed to serve the incoming connection")
		current_session.InboundConn.Close()
		return
	}

	HandleConnection(current_session)
}