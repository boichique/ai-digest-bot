package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"digest_bot/internal/client"
	"digest_bot/internal/config"
	"digest_bot/internal/crontasks"
	"digest_bot/internal/validation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
)

func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "parse config")

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	cl := client.New(cfg.BaseURL)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	cr := cron.New()
	cr.AddFunc("0 0 19 * *", func() { crontasks.SendDigestToAllUsers(context.Background(), bot, cl) }) // update every day at 21:00 msk (0 0 19 * *)
	cr.Start()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// log incoming messages with username
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		// add interesting source to db
		case "new":
			source := strings.TrimLeft(update.Message.Text, "/new ")
			if err := validation.ValidateLink(validation.YoutubeChannelRe, source); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if err := cl.CreateSource(update.Message.Chat.ID, source); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source %q added", source)))

		// output list of sources
		case "list":
			sources, err := cl.GetSourcesList(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if len(sources) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No sources found"))
				break
			}

			var list string
			for _, source := range sources {
				list = list + fmt.Sprintf("\nhttps://www.youtube.com/%s", source) // попробовать сделать через стрингбилдер
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Sources list:\n%s", list)))

		// get new videos on sources
		case "newVideos":
			videos, err := cl.GetNewVideosForUserSources(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if len(videos) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No new videos found"))
				break
			}

			var list string
			for _, video := range videos {
				list = list + fmt.Sprintf("\n%q: https://www.youtube.com/watch?v=%s", video.Title, video.VideoID) // попробовать сделать через стрингбилдер
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Today's new videos:\n%s", list)))

		// get digest of all sources
		case "digest":
			digest, err := cl.GetDigestForUserSource(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if len(digest) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No digest for today"))
				break
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Today's digest:\n%s", digest)))

		// delete source by youtube link
		case "delete":
			source := strings.TrimLeft(update.Message.Text, "/delete ")
			if err := validation.ValidateLink(validation.YoutubeChannelRe, source); err != nil {
				log.Print(source, err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			if err := cl.DeleteSourceByLink(update.Message.Chat.ID, source); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				break
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source %q deleted", source)))
		}
	}
}

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
