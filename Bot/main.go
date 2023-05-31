package main

import (
	"fmt"
	"log"
	"strings"

	"digest_bot/internal/client"
	"digest_bot/internal/config"
	"digest_bot/internal/validation"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

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

			if err := client.CreateSource(update.Message.Chat.ID, source); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source %q added", source)))

		// вывод списка источников
		case "list":
			sources, err := client.GetSourcesList(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			var list string
			for _, source := range sources {
				list = list + "\n" + source // попробовать сделать через стрингбилдер
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source list\n%s", list)))

		// удаление источника по ссылке
		case "delete":
			source := strings.Trim(update.Message.Text, "/delete ")
			if err := validation.ValidateYoutube(source); err != nil {
				log.Print(source, err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if err := client.DeleteSourceByLink(update.Message.Chat.ID, source); err != nil {
				log.Print(err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			log.Print(err)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source %q deleted", source)))
		}
	}
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
