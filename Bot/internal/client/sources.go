package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) CreateSource(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Put(c.path("/api/users/%s", strUserID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func (c *Client) GetUsersList() ([]string, error) {
	resp, err := c.client.R().
		Get(c.path("/api/users"))
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

func (c *Client) GetSourcesList(userID int64) ([]string, error) {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		Get(c.path("/api/users/%s/sources", strUserID))
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

func (c *Client) GetNewVideosForUserSources(userID int64) ([]Video, error) {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		Get(c.path("/api/users/%s", strUserID))
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

func (c *Client) GetDigestForUserSource(userID int64) (string, error) {
	strUserID := strconv.Itoa(int(userID))

	var digest string
	resp, err := c.client.R().
		SetResult(&digest).
		Get(c.path("/api/users/%s/digest", strUserID))

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return digest, err
}

func (c *Client) DeleteSourceByLink(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Delete(c.path("/api/users/%s", strUserID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
