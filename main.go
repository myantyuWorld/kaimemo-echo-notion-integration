package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jomei/notionapi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type KaimemoResponse struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

// export XXは、開いているターミナルのみ有効
// export PATH=$PATH:$(go env GOPATH)/bin && air -c .air.tomlでホットリロードを有効化
func main() {
	e := echo.New()
	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fmt.Println(os.Getenv("FRONTEND_URL"))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", os.Getenv("FRONTEND_URL")},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodDelete},
	}))

	apiKey := os.Getenv("NOTION_API_KEY")
	databaseID := os.Getenv("NOTION_DATABASE_ID")

	client := notionapi.NewClient(notionapi.Token(apiKey))
	query := &notionapi.DatabaseQueryRequest{}

	resp, err := client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), query)
	if err != nil {
		log.Fatalf("failed to notion query database: %v", err)
	}

	var kaimemoResponses []KaimemoResponse
	for _, result := range resp.Results {
		properties := result.Properties

		data := KaimemoResponse{}
		for key, property := range properties {
			switch prop := property.(type) {
			case *notionapi.TitleProperty:
				fmt.Printf("  %s (Title): ", key)
				for _, text := range prop.Title {
					fmt.Print(text.Text.Content)
					data.Name = text.Text.Content
				}
				fmt.Println()
			case *notionapi.RichTextProperty:
				fmt.Printf("  %s (RichText): ", key)
				for _, text := range prop.RichText {
					fmt.Print(text.Text.Content)
				}
				fmt.Println()
			case *notionapi.SelectProperty:
				fmt.Printf("  %s (Select): %s\n", key, prop.Select.Name)
			case *notionapi.MultiSelectProperty:
				fmt.Printf("  %s (MultiSelect): ", key)
				for _, option := range prop.MultiSelect {
					fmt.Printf("%s ", option.Name)
					data.Tag = option.Name
				}
				fmt.Println()
			default:
				fmt.Printf("  %s: Unhandled property type\n", key)
			}
		}
		kaimemoResponses = append(kaimemoResponses, data)
	}

	spew.Dump(kaimemoResponses)

	e.GET("/kaimemo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, kaimemoResponses)
	})

	// Vercel が使用するポートを環境変数から取得
	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
