package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func GetNewVideosForUserSource(userSource string, youtubeApiToken string) ([]Video, error) {
	channelLink := strings.TrimLeft(userSource, "https://www.youtube.com/@")
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/search?part=snippet&q=%s&key=%s", channelLink, youtubeApiToken))
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
		Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?channelId=%s&part=snippet,id&order=date&maxResults=20&key=%s", channelID, youtubeApiToken))
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
