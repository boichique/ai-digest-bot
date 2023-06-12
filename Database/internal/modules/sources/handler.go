package sources

import (
	"log"
	"net/http"

	"digest_bot_database/internal/echox"

	"github.com/labstack/echo/v4"
)

const youtubeVideoURL = "https://www.youtube.com/watch?v="

type Handler struct {
	service *Service
	client  *Client
}

func NewHandler(service *Service, client *Client) *Handler {
	return &Handler{
		service: service,
		client:  client,
	}
}

func (h *Handler) CreateSource(c echo.Context) error {
	req, err := echox.Bind[PutAndDeleteRequest](c)
	if err != nil {
		return err
	}

	source := &Source{
		UserID: req.UserID,
		Source: req.Source,
	}

	return h.service.CreateSource(c.Request().Context(), source)
}

func (h *Handler) GetUsersIDList(c echo.Context) error {
	users, err := h.service.GetUsersIDList(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUserSourcesByID(c echo.Context) error {
	req, err := echox.Bind[GetRequest](c)
	if err != nil {
		return err
	}

	sources, err := h.service.GetSourcesByUserID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) GetDigestForUserSource(c echo.Context) error {
	req, err := echox.Bind[GetRequest](c)
	if err != nil {
		return err
	}

	sourcesIDs, err := h.service.GetSourcesByUserID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	var list []Video
	for _, sourceID := range sourcesIDs {
		videos, err := h.client.GetNewVideosForUserSourceByHour(sourceID)
		if err != nil {
			return err
		}

		list = append(list, videos...)
	}

	var fullDigest string
	for _, video := range list {
		digest, err := GetVideoText(int64(req.UserID), youtubeVideoURL+video.VideoID)
		if err != nil {
			return err
		}

		fullDigest += digest
	}

	digest, err := h.client.GetSourceDigestFromChatGPT(c.Request().Context(), fullDigest)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, digest)
}

func (h *Handler) GetNewVideosForUserSources(c echo.Context) error {
	req, err := echox.Bind[GetRequest](c)
	if err != nil {
		return err
	}

	sources, err := h.service.GetSourcesByUserID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	var list []Video
	for _, source := range sources {
		videos, err := h.client.GetNewVideosForUserSourceByHour(source)
		if err != nil {
			return err
		}
		list = append(list, videos...)

	}

	return c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteSourceByLink(c echo.Context) error {
	req, err := echox.Bind[PutAndDeleteRequest](c)
	if err != nil {
		log.Print(err)
		return err
	}

	source := &Source{
		UserID: req.UserID,
		Source: req.Source,
	}

	return h.service.DeleteSourceByLink(c.Request().Context(), source)
}

type PutAndDeleteRequest struct {
	UserID int    `param:"userID"`
	Source string `json:"source"`
}

type GetRequest struct {
	UserID int `param:"userID"`
}
