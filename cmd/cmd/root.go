package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	"operametrix/mqtt/proxy"
)

type Config struct {
	LocalBroker []proxy.LocalBroker
	Middlewares []string
	Peers []proxy.Peer
    Listeners []proxy.Listener
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "MQTT Proxy for peering",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		var config Config
		viper.Unmarshal(&config)

		for _, peer := range config.Peers {
			go peer.Serve()
		}

		for _, listener := range config.Listeners {
			listener.Serve()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile("conf/mqttproxy.yml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err == nil {
	}
}
