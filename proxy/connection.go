package proxy

import (
	"github.com/spf13/viper"
	"log"
	"net"
	"fmt"
	"operametrix/mqtt/middleware"
	"operametrix/mqtt/session"
	"operametrix/mqtt/notify"
	"github.com/eclipse/paho.mqtt.golang/packets"
)

type LocalBroker struct {
	Hostname string  `yaml:"hostname"`
	Port     int     `yaml:"port"`
}

type LocalBrokerConfig struct {
	Broker LocalBroker
}

func ConnectLocalBroker(current_session *session.Session) (err error) {
	var config LocalBrokerConfig
	viper.Unmarshal(&config)

	localBrokerHost := fmt.Sprintf("%s:%d", config.Broker.Hostname, config.Broker.Port)
	current_session.OutboundConn, err = net.Dial("tcp", localBrokerHost)
	if err != nil {
		log.Println("Failed to contact the broker")
		return err
	}

	return err
}

func HandleConnection(current_session *session.Session) {
	defer current_session.InboundConn.Close()
	defer current_session.OutboundConn.Close()

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
	go SocketReader(current_session.InboundConn, inboundChannel, inboundErrorChannel)

	// Create the routine to manage outbound flow
	outboundChannel := make(chan packets.ControlPacket)
	outboundErrorChannel := make(chan error)
	go SocketReader(current_session.OutboundConn, outboundChannel, outboundErrorChannel)

	go notify.Notify("[" + current_session.CommonName + "] start new session")

	for {
		select {
		case data := <-inboundChannel:
			inboundPipeline.Serve(current_session, &data)

		case <-inboundErrorChannel:
			log.Println("Closed connection from", current_session.InboundConn.RemoteAddr())
			go notify.Notify("[" + current_session.CommonName + "] session error")
			current_session.Destroy()
			return

		case data := <-outboundChannel:
			outboundPipeline.Serve(current_session, &data)

		case <-outboundErrorChannel:
			log.Println("Closed connection from", current_session.OutboundConn.RemoteAddr())
			go notify.Notify("[" + current_session.CommonName + "] session error")
			current_session.Destroy()
			return
		}
	}
}
