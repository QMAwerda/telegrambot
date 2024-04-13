package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/QMAwerda/telegrambot/clients/telegram"
	eventconsumer "github.com/QMAwerda/telegrambot/consumer/event-consumer"
	"github.com/QMAwerda/telegrambot/events/telegram"
	"github.com/QMAwerda/telegrambot/storage/sqlite"
)

// TODO: вывести путь до файлов в параметры конфига
// TODO: передавать tgBotHost как и mustToken из параметров вместо константы
const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	// TODO контекст с таймаутом
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
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
