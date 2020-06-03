package proxy

import (
	"crypto/x509"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
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
	hostString := fmt.Sprintf("%s:%s", listener.Hostname, strconv.Itoa(listener.Port))

	var l net.Listener
	var err error
	if listener.TLS {
		cert, err := tls.LoadX509KeyPair(listener.CertFile, listener.KeyFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if listener.MTLS {
			caCertPool  := x509.NewCertPool()
			caCert, err := ioutil.ReadFile(listener.CAFile)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs  = caCertPool
		}

		log.Println("Start to listen on", hostString, "with TLS")
		l, err = tls.Listen("tcp", hostString, tlsConfig)
	} else {
		log.Println("Start to listen on", hostString)
		l, err = net.Listen("tcp", hostString)
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
		log.Println("[", hostString, "] new connection from", conn.RemoteAddr())
		go LaunchConnection(conn)
	}
}

func LaunchConnection(inboundConn net.Conn) {
	outboundConn, err := ConnectLocalBroker()
	if err != nil {
		log.Println("Failed to serve the incoming connection")
		outboundConn.Close()
		return
	}

	HandleConnection(inboundConn, outboundConn)
}