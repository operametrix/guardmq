package notify

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
)

type TelegramNotify struct {
	Token  string    `yaml:"token"`
	ChatID string    `yaml:"chatid"`
}

func (self *TelegramNotify) Notify(message string) {
	if (len(self.Token) == 0) || (len(self.ChatID) == 0) {
		return
	}

	req := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", self.Token, self.ChatID, message)
	resp, err := http.Get(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
}
