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

func (h *Handler) GetSourceText(c echo.Context) error {
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
