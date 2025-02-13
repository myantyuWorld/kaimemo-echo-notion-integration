package main

import (
	"net/http"
	"template-echo-notion-integration/config"
	"template-echo-notion-integration/internal/handler"
	"template-echo-notion-integration/internal/repository"
	"template-echo-notion-integration/internal/service"

	"github.com/davecgh/go-spew/spew"
	"github.com/jomei/notionapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type KaimemoResponse struct {
	ID   notionapi.ObjectID `json:"id"`
	Tag  string             `json:"tag"`
	Name string             `json:"name"`
	Done bool               `json:"done"`
}

type CreateKaimemoRequest struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

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

	repo := repository.NewNotionRepository(appConfig.NotionAPIKey, appConfig.NotionKaimemoDatabaseID)
	service := service.NewKaimemoService(repo)
	handler := handler.NewKaimemoHandler(service)

	e.GET("/kaimemo", handler.FetchKaimemo)
	e.POST("/kaimemo", handler.CreateKaimemo)
	e.DELETE("/kaimemo/:id", handler.RemoveKaimemo)

	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
