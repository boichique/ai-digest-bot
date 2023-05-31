package echox

import "github.com/labstack/echo/v4"

func Bind[T any](c echo.Context) (*T, error) {
	req := new(T)
	if err := c.Bind(req); err != nil {
		return nil, err
	}

	return req, nil
}
