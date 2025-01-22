package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Vercel!")
	})

	// Vercel が使用するポートを環境変数から取得
	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
