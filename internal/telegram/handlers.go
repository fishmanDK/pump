package telegram

import (
	"fmt"
	"github.com/fishmanDK/pump/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

const (
	commandStart              = "start"
	stateStart                = "start"
	stateDownloadVideoMp4     = "Скачать видео(mp4)"
	stateDownloadVideoMp3     = "Скачать видео(mp3)"
	stateWaitingForLinkForMp4 = "Введите ссылку на видео для формата mp4"
	stateWaitingForLinkForMp3 = "Введите ссылку на видео для формата mp3"
)

type Handlers struct {
	Service service.Service
}

func NewHandlers(service service.Service) *Handlers {
	return &Handlers{
		Service: service,
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	const op = "internal.telegram.bot.handleStartCommand"
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.Start)
	msg.ReplyMarkup = b.Keyboard.GetKeyboard()
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = b.StateMachine.SetupState(message.Chat.ID, stateStart)
	log.Println(err)
	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	const op = "internal.telegram.bot.handleUnknownCommand"
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)
	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	const op = "internal.telegram.bot.handleMessage"

	userState := b.StateMachine.GetState(message.Chat.ID)

	switch userState {
	case stateStart:
		log.Println(message.Text)
		switch message.Text {
		case stateDownloadVideoMp4:
		case stateDownloadVideoMp3:
		default:
			msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)
			_, err := b.bot.Send(msg)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		err := b.StateMachine.SetupState(message.Chat.ID, message.Text)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	userState = b.StateMachine.GetState(message.Chat.ID)

	switch userState {
	case stateStart:
		log.Println(message.Text)
		switch message.Text {
		case stateDownloadVideoMp4:
		case stateDownloadVideoMp3:
		default:
			msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Errors.UnknownCommand)
			_, err := b.bot.Send(msg)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		err := b.StateMachine.SetupState(message.Chat.ID, message.Text)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil

	case stateDownloadVideoMp4:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.InputUrl)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = b.StateMachine.SetupState(msg.ChatID, stateWaitingForLinkForMp4)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	case stateDownloadVideoMp3:
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.InputUrl)
		_, err := b.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = b.StateMachine.SetupState(msg.ChatID, stateWaitingForLinkForMp3)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	case stateWaitingForLinkForMp4:
		file, videoTitle, err := b.service.DownloadMp4(message.Text)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		fileBytes := tgbotapi.FileBytes{
			Name:  "file_name",
			Bytes: make([]byte, 0),
		}
		file, err = os.Open(fmt.Sprintf("videos/%s.mp4", videoTitle)) // открыть файл перед чтением
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSize := fileInfo.Size()
		buffer := make([]byte, fileSize)
		_, err = file.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fileBytes.Bytes = buffer

		err = os.Remove(fmt.Sprintf("videos/%s.mp4", videoTitle))
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		b.StateMachine.SetupState(message.Chat.ID, stateStart)

		video := tgbotapi.NewVideo(message.Chat.ID, fileBytes)
		video.ReplyMarkup = b.Keyboard.GetKeyboard()
		_, err = b.bot.Send(video)
		if err != nil {
			log.Println(fmt.Errorf("%s: %w", op, err))
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	case stateWaitingForLinkForMp3:
		stream, err := b.service.DbownloadMp3(message.Text)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		fileBytes := tgbotapi.FileBytes{
			Name:  "file_name",
			Bytes: stream,
		}

		b.StateMachine.SetupState(message.Chat.ID, stateStart)

		audio := tgbotapi.NewAudio(message.Chat.ID, fileBytes)
		audio.ReplyMarkup = b.Keyboard.GetKeyboard()
		_, err = b.bot.Send(audio)
		if err != nil {
			log.Println(fmt.Errorf("%s: %w", op, err))
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	default:
		return nil
	}
}

//file, err = os.Open(videoPath) // открыть файл перед чтением
//if err != nil {
//	log.Fatal(err)
//}
//defer file.Close()
//
//fileInfo, err := file.Stat()
//if err != nil {
//	log.Fatal(err)
//}
//fileSize := fileInfo.Size()
//buffer := make([]byte, fileSize)
//_, err = stream.Read(buffer)
//if err != nil {
//	log.Fatal(err)
//}
