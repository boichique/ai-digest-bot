package crontasks

import (
	"context"
	"log"
	"strconv"

	"digest_bot/internal/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendDigestToAllUsers(ctx context.Context, bot *tgbotapi.BotAPI, client *client.Client) {
	users, err := client.GetUsersList()
	if err != nil {
		log.Print(err)
		return
	}

	for _, user := range users {
		userID, err := strconv.Atoi(user)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(int64(userID), "Can't get digest for you"))
		}
		bot.Send(tgbotapi.NewMessage(int64(userID), "Test"))
	}
}
