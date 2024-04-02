package main

import (
	"flag"
	"log"

	tgClient "github.com/QMAwerda/telegrambot/clients/telegram"
	"github.com/QMAwerda/telegrambot/consumer/event-consumer"
	"github.com/QMAwerda/telegrambot/events/telegram"
	"github.com/QMAwerda/telegrambot/storage/files"
)

// TODO: вывести путь до файлов в параметры конфига
// TODO: передавать tgBotHost как и mustToken из параметров вместо константы
const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
