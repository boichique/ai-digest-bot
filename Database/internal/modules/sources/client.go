package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"digest_bot_database/internal/apperrors"
	"digest_bot_database/internal/log"
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
)

const (
	transcriptorURL     = "http://transcriptor:10001/transcribe"
	youtubeAPIsearchURL = "https://youtube.googleapis.com/youtube/v3/search?"
)

func GetDigestFromChatGPT(ctx context.Context, fullDigest string, chatGPTApiToken string) (string, error) {
	log.FromContext(ctx).Info(
		"digest for chatGPT",
		"digest", fullDigest,
	)

	client := openai.NewClient(chatGPTApiToken)
	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			Prompt:      fmt.Sprintf("ОЧЕНЬ КРАТКО ПЕРЕСКАЖИ ТЕКСТ В 100 СИМВОЛОВ: \n\n%s", fullDigest),
			Temperature: 1,
			MaxTokens:   300,
		},
	)
	if err != nil {
		log.FromContext(ctx).Error(
			"chatGPT response error: ",
			"error", err.Error(),
		)
		return "", fmt.Errorf("ChatCompletion error: %w", err)
	}

	return resp.Choices[0].Text, nil
}

func GetNewVideosForUserSource(sourceID string, youtubeApiToken string) ([]Video, error) {
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("%spart=snippet&type=channel&q=%s&key=%s", youtubeAPIsearchURL, sourceID, youtubeApiToken))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, apperrors.Internal(err)
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
		Get(fmt.Sprintf("%spart=snippet,id&channelId=%s&order=date&maxResults=15&key=%s", youtubeAPIsearchURL, channelID, youtubeApiToken))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var searchListResponse SearchListResponse
	err = json.Unmarshal(resp.Body(), &searchListResponse)
	if err != nil {
		return nil, err
	}

	var videos []Video
	today := time.Now().Add(-24 * time.Hour)
	for _, item := range searchListResponse.Items {
		video := Video{
			Title:       item.Snippet.Title,
			VideoID:     item.ID.VideoID,
			PublishedAt: item.Snippet.PublishedAt,
		}

		date, err := time.Parse(time.RFC3339, video.PublishedAt)
		if err != nil {
			return nil, apperrors.Internal(err)
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
		Post(transcriptorURL)
	if err != nil {
		return "", apperrors.Internal(err)
	}

	return resp.String(), nil
}
