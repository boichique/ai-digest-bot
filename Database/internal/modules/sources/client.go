package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	openai "github.com/sashabaranov/go-openai"
)

func GetNewVideosForUserSource(userSource string, youtubeApiToken string) ([]Video, error) {
	channelLink := strings.TrimLeft(userSource, "https://www.youtube.com/@")
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/search?part=snippet&type=channel&q=%s&key=%s", channelLink, youtubeApiToken))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var data map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}

	var channelID string
	for _, item := range data["items"].([]interface{}) {
		itemMap := item.(map[string]interface{})
		if itemMap["id"].(map[string]interface{})["kind"] == "youtube#channel" {
			channelID = itemMap["id"].(map[string]interface{})["channelId"].(string)
			break
		}
	}

	resp, err = client.R().
		Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?channelId=%s&part=snippet,id&order=date&maxResults=15&key=%s", channelID, youtubeApiToken))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var searchListResponse SearchListResponse
	err = json.Unmarshal(resp.Body(), &searchListResponse)
	if err != nil {
		return nil, err
	}

	today := time.Now().Add(-24 * time.Hour)
	var videos []Video
	for _, item := range searchListResponse.Items {
		video := Video{
			Title:       item.Snippet.Title,
			VideoID:     item.Id.VideoId,
			PublishedAt: item.Snippet.PublishedAt,
		}

		date, err := time.Parse(time.RFC3339, video.PublishedAt)
		if err != nil {
			return nil, err
		}

		if date.After(today) {
			videos = append(videos, video)
		}
	}

	return videos, nil
}

func GetVideoText(userID int64, source string) (string, error) {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"link":   source,
			"output": strUserID,
		}).
		Post("http://transcriptor:10001/transcribe")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error getting text from video: %s", resp.Status())
	}

	return resp.String(), nil
}

func GetDigestFromChatGPT(fullDigest string, chatGPTApiToken string) (string, error) {
	query := "Summarize this text in 200-300 symbols: "
	log.Print(fullDigest)

	client := openai.NewClient(chatGPTApiToken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query + fullDigest,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %w\n", err)
	}

	return resp.Choices[0].Message.Content, nil
}
