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

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "GuardMQ - MQTT Proxy for peering v1.0.0",
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/etc/guardmq/guardmq.yml", "config file")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	viper.ReadInConfig();
}
