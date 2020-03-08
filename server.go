package main_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
)

func TestMain(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
