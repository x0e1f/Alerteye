package telegram

/*
The MIT License (MIT)

Copyright (c) 2020 Davide Pataracchia

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/x0e1f/Alerteye/configs"
	"github.com/x0e1f/Alerteye/database"
	"log"
	"strconv"
	"time"
)

// StartConsumer :: Start Alerteye Telegram consumer
func StartConsumer(dbPath string, configs *configs.Config) {
	log.Print("Telegram consumer started")
	ChatID, _ := strconv.ParseInt(configs.ChatID, 10, 64)

	for {
		bot, err := tgbotapi.NewBotAPI(configs.BotToken)
		if err == nil {
			alert, _ := database.AlertToSend(dbPath)
			if alert.URL != "" {
				messageBody := "[" + alert.Source + "]\n"
				if alert.Topic != "" {
					messageBody += "[" + alert.Topic + "]\n"
				}
				messageBody += "[" + alert.Title + "]\n"
				messageBody += alert.URL
				msg := tgbotapi.NewMessage(ChatID, messageBody)
				bot.Send(msg)
				database.AlertSent(dbPath, alert.URL)
				log.Print("[Telegram] " + alert.Title)
			}
		} else {
			log.Print("Telegram error: ", err)
		}

		time.Sleep(time.Duration(configs.SendTime) * time.Minute)
	}
}
