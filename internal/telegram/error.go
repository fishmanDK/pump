package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) commandError(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.messages.Errors.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) messageError(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.messages.Errors.BadUrl)
	_, err := b.bot.Send(msg)

	return err
}