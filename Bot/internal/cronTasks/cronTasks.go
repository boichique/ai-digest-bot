package cronTasks

import (
	"log"
	"strconv"

	"digest_bot/internal/client"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Test(bot *tgbotapi.BotAPI) error {
	log.Print("Считываю пользователей с базы и делаю рассылку")
	usersID, err := client.GetUsersList()
	if err != nil {
		log.Print(err)
		return err
	}

	for _, userID := range usersID {
		strUserID, err := strconv.Atoi(userID)
		if err != nil {
			log.Print(err)
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(int64(strUserID), "Вы в списках.")); err != nil {
			log.Print("Ошибка при отправке сообщения")
		}
	}

	return nil
}
