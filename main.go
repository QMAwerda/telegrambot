package main

import (
	"flag"
	"log"

	"github.com/QMAwerda/telegrambot/clients/telegram"
)

// Стоит передавать хост как и mustToken из параметров, так приложение будет более гибким
const (
	tgBotHost = "api.telegram.org"
)

func main() {

	//token = flags.Get(token)

	//tgClient = telegram.New(token)

	//fetcher = fetcher.New()
	//processor = processor.New()

	//consumer.Start(fetcher, processor)

	// Реализуем:

	tgClient := telegram.New(tgBotHost, mustToken()) // создае тг клиента

}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
