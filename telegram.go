package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineQueryButton struct {
	Text   string              `json:"text"`
	WebApp tgbotapi.WebAppInfo `json:"web_app"`
}

type InlineQueryAnswerWithButton struct {
	tgbotapi.InlineConfig
	Button InlineQueryButton `json:"button"`
}

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
		if update.InlineQuery != nil {
			params := make(tgbotapi.Params)

			params["inline_query_id"] = update.InlineQuery.ID
			params.AddNonZero("cache_time", 0)
			params.AddInterface("button", InlineQueryButton{
				Text:   "Лист персонажа",
				WebApp: tgbotapi.WebAppInfo{"http://127.0.0.1:8080/sheet"},
			})
			params.AddInterface("results", []interface{}{
				tgbotapi.NewInlineQueryResultCachedPhoto("1", "AgACAgIAAxkBAAMcZrpqn0hJ52L6JOxRTssoMdz_fXwAAqmnMRtOAAHRSd-rUD-IFhl6AQADAgADbQADNQQ"),
				tgbotapi.NewInlineQueryResultCachedPhoto("2", "AgACAgIAAxkBAAMcZrpqn0hJ52L6JOxRTssoMdz_fXwAAqmnMRtOAAHRSd-rUD-IFhl6AQADAgADbQADNQQ"),
				tgbotapi.NewInlineQueryResultCachedPhoto("3", "AgACAgIAAxkBAAMcZrpqn0hJ52L6JOxRTssoMdz_fXwAAqmnMRtOAAHRSd-rUD-IFhl6AQADAgADbQADNQQ"),
				tgbotapi.NewInlineQueryResultCachedPhoto("4", "AgACAgIAAxkBAAMcZrpqn0hJ52L6JOxRTssoMdz_fXwAAqmnMRtOAAHRSd-rUD-IFhl6AQADAgADbQADNQQ"),
				tgbotapi.NewInlineQueryResultCachedPhoto("5", "AgACAgIAAxkBAAMcZrpqn0hJ52L6JOxRTssoMdz_fXwAAqmnMRtOAAHRSd-rUD-IFhl6AQADAgADbQADNQQ"),
			})

			if _, err := bot.MakeRequest("answerInlineQuery", params); err != nil {
				log.Println(err)
			}
			continue
		}
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		}
	}
}
