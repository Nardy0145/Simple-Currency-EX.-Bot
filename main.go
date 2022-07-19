package main

import (
	"encoding/json"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
)

func main() {
	var tokenFlag = flag.String("token", "0", "Bot API token")
	flag.Parse()
	if *tokenFlag == "0" {
		log.Panic("Invalid flag.")
	}
	bot, err := tgbotapi.NewBotAPI(*tokenFlag)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Successfuly authorized into %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ParseMode = "HTML"
		switch update.Message.Command() {
		case "start":
			msg.Text = "<b>Hello!</b>" +
				"\nThis is a bot, specified for currency conversion." +
				"\n<i>You could see an example by using /help</i>"
			bot.Send(msg)
		case "help":
			msg.Text = "To see live time currency exchange rate, you need to use <code>/convert</code>\n" +
				"For example, if you want to see <b>€ to $</b> ex-rate, just type \n" +
				"<code>/convert EUR USD</code>"
			bot.Send(msg)

		case "convert":
			var args []string = strings.Split(update.Message.Text, " ")
			if len(args) == 1 {
				msg.Text = "<b>Looks like you forgot to specify currency</b>"
				bot.Send(msg)
				continue
			}
			if len(args) == 2 {
				msg.Text = "<b>Looks like you forgot to specify the second currency</b>"
				bot.Send(msg)
				continue
			}
			msg.Text = "the current " + args[1] + "/" + args[2] + " ratio is: <b>" + convertCurrency(args[1], args[2]) + "</b>"
			bot.Send(msg)
		}

	}
}

type parsedResp struct {
	Result float64 `json:"result"`
}

func convertCurrency(from, to string) string {
	client := http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.exchangerate.host/convert?from=%s&to=%s", from, to), nil)
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	var result parsedResp
	json.NewDecoder(response.Body).Decode(&result)
	return fmt.Sprintf("%.2f", result.Result)
}
