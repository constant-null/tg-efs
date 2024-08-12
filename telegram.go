package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const TestAPIEndpoint = "https://api.telegram.org/bot%s/test/%s"

func runTGBot() {
	var bot *tgbotapi.BotAPI
	var err error
	if os.Getenv("DEBUG") != "" {
		bot, err = tgbotapi.NewBotAPIWithClient(os.Getenv("TELEGRAM_BOT_KEY"), TestAPIEndpoint, &http.Client{})
	} else {
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_KEY"))
	}
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           update.Message.Chat.ID,
					ReplyToMessageID: update.Message.MessageID,
				},
				Text:                  update.Message.Text,
				DisableWebPagePreview: false,
				ParseMode:             "HTML",
			}
			bot.Send(msg)
		}
	}
}
