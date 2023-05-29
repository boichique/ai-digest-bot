package main

import (
	"digest_bot_database/internal/config"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"digest_bot/internal/validation"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	youtubeRe = `^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$`
)

var telegramBotToken string

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
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

			err = createSource(strconv.Itoa(int(update.Message.Chat.ID)), source)
			if err != nil {
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

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func createSource(userID string, source string) error {
	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Put(fmt.Sprintf("http://localhost:8000/api/%s", userID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
