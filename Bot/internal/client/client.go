package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const baseURL = "http://server:10000/api/users/"

func CreateSource(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Put(baseURL + strUserID)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func GetUsersList() ([]string, error) {
	client := resty.New()
	resp, err := client.R().
		Get(baseURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var users []string
	err = json.Unmarshal(resp.Body(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetSourcesList(userID int64) ([]string, error) {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("%s%s/sources", baseURL, strUserID))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var sources []string
	err = json.Unmarshal(resp.Body(), &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}

// delete this after tests
func GetVideoText(userId int64, source string) (string, error) {
	strUserID := strconv.Itoa(int(userId))

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

	log.Print("message received")
	log.Print(resp.String())
	return resp.String(), nil
}

func GetNewVideosForUserSources(userID int64) ([]Video, error) {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		Get(baseURL + strUserID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var videos []Video
	err = json.Unmarshal(resp.Body(), &videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func DeleteSourceByLink(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Delete(baseURL + strUserID)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
