package main

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func main() {
	channelID := "CNN"
	fmt.Println(channelID)
	// Создаем клиент YouTube API
	client := &http.Client{
		Transport: &transport.APIKey{
			Key: "AIzaSyCgdt_WOaH_1SmLI5-9WmpJSeEY4mTTsRk",
		},
	}

	// Создаем объект сервиса YouTube API
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Выполняем запрос на получение списка видео на канале
	call := service.Search.List([]string{"id", "snippet"}).
		ChannelId("CHANNEL_ID").
		MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Выводим список видео в консоль
	for _, item := range response.Items {
		fmt.Println(item.Snippet.Title)
	}
}
