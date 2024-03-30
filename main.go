package main

//TODO: Стоит передавать host как и mustToken из параметров вместо const tgBotHost

import (
	"flag"
	"log"

	"github.com/QMAwerda/telegrambot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {

	//token = flags.Get(token) - done

	//tgClient = telegram.New(token) - done

	//fetcher = fetcher.New()
	//processor = processor.New()

	//consumer.Start(fetcher, processor)

	// Реализуем:

	tgClient := telegram.New(tgBotHost, mustToken())

}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
