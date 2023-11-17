package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type keyboards tgbotapi.ReplyKeyboardMarkup

type Keyboard interface {
	GetKeyboard() tgbotapi.ReplyKeyboardMarkup
}

func NewKeyboard() *keyboards { return &keyboards{} }

func (k keyboards) GetKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Скачать видео(mp4)"),
			tgbotapi.NewKeyboardButton("Скачать видео(mp3)"),
		),
	)

	return keyboard
}
