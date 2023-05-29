package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"digest_bot/internal/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	youtubeRe = `^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$`
)

var telegramBotToken string

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		// добавление интересующего источника
		case "new":
			source := strings.Trim(update.Message.Text, "/new ")
			if err := validation.ValidateYoutube(source); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source %q added", source)))

		// вывод дайджеста
		case "digest":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Not implemented"))
		}

	}
}

func init() {
	flag.StringVar(&telegramBotToken, "t", "", "Telegram Bot Token")
	flag.Parse()

	if telegramBotToken == "" {
		log.Fatal("telegram bot token is required")
	}
}
