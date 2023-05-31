package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func CreateSource(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Put(fmt.Sprintf("http://server:10000/api/%s", strUserID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func GetSourcesList(userID int64) ([]string, error) {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("http://server:10000/api/%s/digest", strUserID))
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

func DeleteSourceByLink(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	client := resty.New()
	resp, err := client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Delete(fmt.Sprintf("http://server:10000/api/%s", strUserID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
