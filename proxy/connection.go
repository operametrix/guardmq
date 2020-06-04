package proxy

import (
	"github.com/spf13/viper"
	"log"
	"net"
	"fmt"
	"operametrix/mqtt/middleware"
	"operametrix/mqtt/session"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

type LocalBroker struct {
	Hostname string  `yaml:"hostname"`
	Port     int     `yaml:"port"`
}

type LocalBrokerConfig struct {
	Broker LocalBroker
}

func ConnectLocalBroker() (conn net.Conn, err error) {
	var config LocalBrokerConfig
	viper.Unmarshal(&config)

	localBrokerHost := fmt.Sprintf("%s:%d", config.Broker.Hostname, config.Broker.Port)
	outboundConn, err := net.Dial("tcp", localBrokerHost)
	if err != nil {
		log.Println("Failed to contact the broker")
		return
	}
	
	return outboundConn, err
}

func HandleConnection(inboundConn net.Conn, outboundConn net.Conn) {
	defer inboundConn.Close()
	defer outboundConn.Close()

	var current_session session.Session
	session.ActiveSessions = append(session.ActiveSessions, &current_session)
	current_session.InboundConn = inboundConn
	current_session.OutboundConn = outboundConn

	var config middleware.MiddlewareConfig
	viper.Unmarshal(&config)

	// Create the chain of middleware for inbound
	var inboundPipeline middleware.Handler
	inboundPipeline = middleware.HandlerFunc(ForwardToBroker)
	for _, m := range config.Middlewares {
		switch m {
		case "LoggingMiddleware":
			inboundPipeline = middleware.LoggingMiddleware(inboundPipeline)
		case "ActiveMiddleware":
			inboundPipeline = middleware.ActiveMiddleware(inboundPipeline)
		case "ExampleMiddleware":
			inboundPipeline = middleware.ExampleMiddleware(inboundPipeline)
		case "StatsMiddleware":
			inboundPipeline = middleware.StatsInMiddleware(inboundPipeline)
		}
	}

	// Create the chain of middleware for outbound
	var outboundPipeline middleware.Handler
	outboundPipeline = middleware.HandlerFunc(ForwardToClient)
	for _, m := range config.Middlewares {
		switch m {
		case "LoggingMiddleware":
			outboundPipeline = middleware.LoggingMiddleware(outboundPipeline)
		case "ActiveMiddleware":
			outboundPipeline = middleware.ActiveMiddleware(outboundPipeline)
		case "ExampleMiddleware":
			outboundPipeline = middleware.ExampleMiddleware(outboundPipeline)
		case "StatsMiddleware":
			outboundPipeline = middleware.StatsOutMiddleware(outboundPipeline)
		}
	}

	// Create the routine to manage inbound flow
	inboundChannel := make(chan packets.ControlPacket)
	inboundErrorChannel := make(chan error)
	go SocketReader(inboundConn, inboundChannel, inboundErrorChannel)

	// Create the routine to manage outbound flow
	outboundChannel := make(chan packets.ControlPacket)
	outboundErrorChannel := make(chan error)
	go SocketReader(outboundConn, outboundChannel, outboundErrorChannel)

	for {
		select {
		case data := <-inboundChannel:
			inboundPipeline.Serve(&current_session, &data)

		case <-inboundErrorChannel:
			current_session.Destroy()
			log.Println("Closed connection from", inboundConn.RemoteAddr())
			return

		case data := <-outboundChannel:
			outboundPipeline.Serve(&current_session, &data)

		case <-outboundErrorChannel:
			current_session.Destroy()
			log.Println("Closed connection from", inboundConn.RemoteAddr())
			return
		}
	}
}
