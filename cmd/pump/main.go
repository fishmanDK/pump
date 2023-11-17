package main

import (
	"github.com/fishmanDK/pump/configs"
	service2 "github.com/fishmanDK/pump/internal/service"
	"github.com/fishmanDK/pump/internal/state_machine"
	"github.com/fishmanDK/pump/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type OptsBot struct {
}

func main() {
	cfg, err := configs.MustConfig()

	// TODO: init logger

	if err != nil {
		log.Fatal(err)
	}
	redis, err := state_machine.NewStateMachine()

	service := service2.NewService()

	botApi, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err.Error())
	}

	bot := telegram.NewBot(botApi, nil, cfg.Messages, tgbotapi.ReplyKeyboardMarkup{}, service, redis)

	err = bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
