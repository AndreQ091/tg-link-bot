package main

import (
	"flag"
	"log"
	tgClient "tg-link-bot/clients/telegram"
	eventConsumer "tg-link-bot/consumer/event-consumer"
	"tg-link-bot/events/telegram"
	"tg-link-bot/storage/files"
)

const (
	tgHost      = "api.telegram.org"
	batchSize   = 100
	storagePath = "files_storage"
)

func main() {
	tgClient := tgClient.New(tgHost, mustToken())
	eventProcessor := telegram.NewProcessor(tgClient, files.New(storagePath))

	log.Println("server started")

	if err := eventConsumer.New(eventProcessor, eventProcessor, batchSize).Start(); err != nil {
		log.Fatal("fatal error", err)
	}
}

func mustToken() string {
	token := flag.String("bot-token", "", "token for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("telegram token is absent")
	}

	return *token
}
