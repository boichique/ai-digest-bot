package crontasks

import (
	"strconv"

	"digest_bot/internal/client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendDigestToAllUsers(bot *tgbotapi.BotAPI, client *client.Client) error {
	users, err := client.GetUsersList()
	if err != nil {
		return err
	}

	for _, user := range users {
		userID, err := strconv.Atoi(user)
		if err != nil {
			continue
		}
		bot.Send(tgbotapi.NewMessage(int64(userID), "Test"))
	}

	return nil
}
