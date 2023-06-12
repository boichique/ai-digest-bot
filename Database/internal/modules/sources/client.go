package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"digest_bot_database/internal/apperrors"
	"digest_bot_database/internal/config"
	"digest_bot_database/internal/log"

	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	Config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{Config: cfg}
}

const (
	transcriptorURL                 = "http://transcriptor:10001/transcribe"
	youtubeAPIsearchURL             = "https://youtube.googleapis.com/youtube/v3/search?"
	youtubeAPIgetChannelIDURL       = youtubeAPIsearchURL + "part=snippet&type=channel&"
	youtubeAPIgetNewVideosByHourURL = youtubeAPIsearchURL + "part=snippet,id&order=date&maxResults=15&"
)

func (c *Client) GetSourceDigestFromChatGPT(ctx context.Context, fullDigest string) (string, error) {
	log.FromContext(ctx).Info(
		"digest for chatGPT",
		"digest", fullDigest,
	)

	client := openai.NewClient(c.Config.ChatGPTApiToken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: "Wait for the full text and summarize it in one sentence. End of text will be when I write 'End of text'.",
				},
				{
					Role:    "user",
					Content: fullDigest,
				},
				{
					Role:    "user",
					Content: "End of text",
				},
			},
			Temperature: 0,
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

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) GetNewVideosForUserSourceByHour(sourceID string) ([]Video, error) {
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("%sq=%s&key=%s", youtubeAPIgetNewVideosByHourURL, sourceID, c.Config.YoutubeApiToken))
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
		Get(fmt.Sprintf("%schannelId=%s&key=%s", youtubeAPIsearchURL, channelID, c.Config.YoutubeApiToken))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var searchListResponse SearchListResponse
	err = json.Unmarshal(resp.Body(), &searchListResponse)
	if err != nil {
		return nil, err
	}

	var videos []Video
	today := time.Now().Add(-1 * time.Hour)
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
