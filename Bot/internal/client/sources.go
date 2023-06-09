package client

import (
	"encoding/json"
	"strconv"
)

func (c *Client) CreateSource(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	_, err := c.client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Put(c.path("/api/users/%s", strUserID))

	return err
}

func (c *Client) GetUsersList() ([]string, error) {
	resp, err := c.client.R().
		Get(c.path("/api/users"))
	if err != nil {
		return nil, err
	}

	var users []string
	err = json.Unmarshal(resp.Body(), &users)

	return users, err
}

func (c *Client) GetSourcesList(userID int64) ([]string, error) {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		Get(c.path("/api/users/%s/sources", strUserID))
	if err != nil {
		return nil, err
	}

	var sources []string
	err = json.Unmarshal(resp.Body(), &sources)

	return sources, err
}

func (c *Client) GetNewVideosForUserSources(userID int64) ([]Video, error) {
	strUserID := strconv.Itoa(int(userID))

	resp, err := c.client.R().
		Get(c.path("/api/users/%s", strUserID))
	if err != nil {
		return nil, err
	}

	var videos []Video
	err = json.Unmarshal(resp.Body(), &videos)

	return videos, err
}

func (c *Client) GetDigestForUserSource(userID int64) (string, error) {
	strUserID := strconv.Itoa(int(userID))

	var digest string
	_, err := c.client.R().
		SetResult(&digest).
		Get(c.path("/api/users/%s/digest", strUserID))

	return digest, err
}

func (c *Client) DeleteSourceByLink(userID int64, source string) error {
	strUserID := strconv.Itoa(int(userID))

	_, err := c.client.R().
		SetBody(map[string]string{
			"source": source,
		}).
		Delete(c.path("/api/users/%s", strUserID))

	return err
}
