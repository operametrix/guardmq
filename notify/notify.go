package notify

import (
	"github.com/spf13/viper"
	"log"
)

type NotifyPlugins struct {
	Telegram TelegramNotify
}

type NotifyConfig struct {
	Notify NotifyPlugins
}

func Notify(message string) {
	var config NotifyConfig
	viper.Unmarshal(&config)

	config.Notify.Telegram.Notify(message)
}