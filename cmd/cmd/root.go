package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"

	"operametrix/mqtt/proxy"
	"operametrix/mqtt/middleware"
)

type Config struct {
	LocalBroker []proxy.LocalBroker
	Middlewares []string
	Peers []proxy.Peer
    Listeners []proxy.Listener
}

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "guardmq",
	Short: "GuardMQ - MQTT Proxy toolbox v1.0.0",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		var config Config
		viper.Unmarshal(&config)

		for _, peer := range config.Peers {
			go peer.Serve()
		}

		for _, listener := range config.Listeners {
			go listener.Serve()
		}

		if len(config.Peers) == 0 && len(config.Listeners) == 0 {
			log.Println("Error: no peer and no listener defined")
			os.Exit(1)
		}

		// Start Middleware routines
		for _, m := range config.Middlewares {
			switch m {
			case "StatsMiddleware":
				go middleware.StatsRoutineMiddleware()
			}
		}

		// Wait SIGTERM signal
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan

		log.Println("Closed the proxy")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/etc/guardmq/guardmq.yml", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	viper.ReadInConfig();
}
