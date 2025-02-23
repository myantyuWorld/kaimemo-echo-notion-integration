package main

import (
	"net/http"
	"template-echo-notion-integration/config"
	"template-echo-notion-integration/internal/handler"
	"template-echo-notion-integration/internal/repository"
	"template-echo-notion-integration/internal/service"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// export XXは、開いているターミナルのみ有効
// export PATH=$PATH:$(go env GOPATH)/bin && air -c .air.toml でホットリロードを有効化
func main() {
	appConfig := config.LoadConfig()
	spew.Dump(appConfig)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConfig.AllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodDelete},
	}))

	repo := repository.NewNotionRepository(
		appConfig.NotionAPIKey,
		appConfig.NotionKaimemoDatabaseInputID,
		appConfig.NotionKaimemoDatabaseSummaryRecordID,
	)
	service := service.NewKaimemoService(repo)
	handler := handler.NewKaimemoHandler(service)

	e.GET("/kaimemo", handler.FetchKaimemo)
	e.POST("/kaimemo", handler.CreateKaimemo)
	e.DELETE("/kaimemo/:id", handler.RemoveKaimemo)

	e.GET("/ws/kaimemo", handler.WebsocketTelegraph)

	e.GET("/kaimemo/summary", handler.FetchKaimemoSummaryRecord)
	e.POST("/kaimemo/summary", handler.CreateKaimemoAmount)
	e.DELETE("/kaimemo/summary/:id", handler.RemoveKaimemoAmount)

	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
