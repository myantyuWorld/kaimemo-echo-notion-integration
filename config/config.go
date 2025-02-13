package config

import (
	"log"
	"os"
)

type AppConfig struct {
	Port                    string
	NotionAPIKey            string
	NotionKaimemoDatabaseID string
	AllowOrigins            []string
}

func LoadConfig() *AppConfig {
	port := "3000"

	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY is not set")
	}

	notionKaimemoDatabaseID := os.Getenv("NOTION_DATABASE_ID")
	if notionKaimemoDatabaseID == "" {
		log.Fatal("NOTION_DATABASE_ID is not set")
	}

	frontEndUrl := os.Getenv("FRONTEND_URL")
	if frontEndUrl == "" {
		log.Fatal("FRONTEND_URL is not set")
	}

	return &AppConfig{
		Port:                    port,
		NotionAPIKey:            apiKey,
		NotionKaimemoDatabaseID: notionKaimemoDatabaseID,
		AllowOrigins: []string{
			"http://localhost:5173", "http://localhost:4173", frontEndUrl,
		},
	}
}
