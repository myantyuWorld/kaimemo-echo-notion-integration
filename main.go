package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

// export XXは、開いているターミナルのみ有効
// export PATH=$PATH:$(go env GOPATH)/bin && air -c .air.tomlでホットリロードを有効化
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", os.Getenv("FRONTEND_URL")},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodDelete},
	}))

	apiKey := os.Getenv("NOTION_API_KEY")
	databaseID := os.Getenv("NOTION_DATABASE_ID")

	client := notionapi.NewClient(notionapi.Token(apiKey))
	query := &notionapi.DatabaseQueryRequest{}

	e.GET("/kaimemo", func(c echo.Context) error {
		resp, err := client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), query)
		if err != nil {
			log.Fatalf("failed to notion query database: %v", err)
		}

		var kaimemoResponses []KaimemoResponse
		for _, result := range resp.Results {
			properties := result.Properties

			data := KaimemoResponse{}
			data.ID = result.ID
			for _, property := range properties {
				switch prop := property.(type) {
				case *notionapi.TitleProperty:
					for _, text := range prop.Title {
						data.Name = text.Text.Content
					}
				case *notionapi.SelectProperty:
					data.Tag = prop.Select.Name
				case *notionapi.CheckboxProperty:
					data.Done = prop.Checkbox
				default:
					// fmt.Printf("  %s: Unhandled property type\n", key)
				}
			}
			kaimemoResponses = append(kaimemoResponses, data)
		}

		return c.JSON(http.StatusOK, kaimemoResponses)
	})

	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
