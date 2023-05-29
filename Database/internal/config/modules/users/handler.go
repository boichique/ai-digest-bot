package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) GetSourcesByID(c echo.Context) error {
	var req UserIDRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.service.GetUserSourcesByID(c.Request().Context(), req.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateSource(c echo.Context) error {
	var req PutRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &User{
		ID:     req.UserID,
		Source: req.Source,
	}

	return h.service.CreateSource(c.Request().Context(), user)
}

type UserIDRequest struct {
	UserID int `param:"userID"`
}

type PutRequest struct {
	UserID int    `param:"userID"`
	Source string `json:"source"`
}
