package sources

import (
	"log"
	"net/http"

	"digest_bot_database/internal/echox"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSource(c echo.Context) error {
	req, err := echox.Bind[PutAndDeleteRequest](c)
	if err != nil {
		log.Print(err)
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
		log.Print(err)
		return err
	}

	sources, err := h.service.GetUserSourcesByID(c.Request().Context(), req.UserID)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) GetNewVideosForUserSources(c echo.Context) error {
	req, err := echox.Bind[GetRequest](c)
	if err != nil {
		return err
	}

	sources, err := h.service.GetUserSourcesByID(c.Request().Context(), req.UserID)
	if err != nil {
		return err
	}

	var list []Video
	youtubeApiToken := c.Get("YoutubeApiToken").(string)
	for _, source := range sources {
		videos, err := GetNewVideosForUserSource(source, youtubeApiToken)
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
