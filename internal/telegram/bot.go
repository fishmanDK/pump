package telegram

import (
	"fmt"
	"github.com/fishmanDK/pump/configs"
	"github.com/fishmanDK/pump/internal/service"
	"github.com/fishmanDK/pump/internal/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kkdai/youtube/v2"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	messages configs.Messages
	//keyboard tgbotapi.ReplyKeyboardMarkup
	service *service.Service

	Keyboard
	StateMachine  *state_machine.StateMachine
	YoutubeClient youtube.Client
	Logger        *zap.Logger
}

func NewBot(bot *tgbotapi.BotAPI, logger *zap.Logger, messages configs.Messages, keyboard tgbotapi.ReplyKeyboardMarkup, service *service.Service, machine *state_machine.StateMachine) *Bot {
	return &Bot{
		service:  service,
		bot:      bot,
		messages: messages,
		//keyboard: keyboard,

		Keyboard:      NewKeyboard(),
		StateMachine:  machine,
		YoutubeClient: youtube.Client{},
		Logger:        logger,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/") {
			if err := b.handleCommand(update.Message); err != nil {
				b.commandError(update.Message.Chat.ID)
			}

			continue
		}
		//b.Logger.Info("message")
		if err := b.handleMessage(update.Message); err != nil {
			log.Println(err)
			b.messageError(update.Message.Chat.ID)
		} else {

			id := strconv.FormatInt(update.Message.Chat.ID, 10)
			str := fmt.Sprintf("%s: %s", id, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(1000574498, str)
			b.bot.Send(msg)
		}
	}

	return nil
}
