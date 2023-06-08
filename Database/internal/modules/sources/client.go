package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"digest_bot_database/internal/log"
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client  *resty.Client
	baseURL string
}

func NewClient(url string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)
	rc.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if resp.StatusCode() >= 400 {
			return fmt.Errorf("http error %d: %s", resp.StatusCode(), resp.Status())
		}
		return nil
	})

	return &Client{
		client:  rc,
		baseURL: url,
	}
}

const (
	transcriptorURL     = "http://transcriptor:10001/transcribe"
	youtubeAPIsearchURL = "https://youtube.googleapis.com/youtube/v3/search?"
)

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.baseURL+f, args...)
}

func (c *Client) GetDigestFromChatGPT(ctx context.Context, fullDigest string, chatGPTApiToken string) (string, error) {
	query := "Summarize this text in 200-300 symbols: "
	log.FromContext(ctx).Info(
		"chatGPT query",
		"fullDigest", fullDigest,
	)

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

func (c *Client) GetNewVideosForUserSource(sourceID string, youtubeApiToken string) ([]Video, error) {
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("%spart=snippet&type=channel&q=%s&key=%s", youtubeAPIsearchURL, sourceID, youtubeApiToken))
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

	resp, err = c.client.R().
		Get(fmt.Sprintf("%spart=snippet,id&channelId=%s&order=date&maxResults=15&key=%s", youtubeAPIsearchURL, channelID, youtubeApiToken))
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
			return nil, err
		}

		if date.After(today) {
			videos = append(videos, video)
		}
	}

	return videos, nil
}

func (c *Client) GetVideoText(userID int64, source string) (string, error) {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		SetBody(map[string]string{
			"link":   source,
			"output": strUserID,
		}).
		Post(transcriptorURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error getting text from video: %s", resp.Status())
	}

	return resp.String(), nil
}
