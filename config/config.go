package config

import (
	"log"
	"os"
)

type AppConfig struct {
	Port                                 string
	NotionAPIKey                         string
	NotionKaimemoDatabaseInputID         string
	NotionKaimemoDatabaseSummaryRecordID string
	AllowOrigins                         []string
}

func LoadConfig() *AppConfig {
	port := "3000"

	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY is not set")
	}

	notionKaimemoDatabaseInputID := os.Getenv("NOTION_DATABASE_KAIMEMO_INPUT")
	if notionKaimemoDatabaseInputID == "" {
		log.Fatal("NOTION_DATABASE_KAIMEMO_INPUT is not set")
	}

	notionKaimemoDatabaseSummaryRecordID := os.Getenv("NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD")
	if notionKaimemoDatabaseSummaryRecordID == "" {
		log.Fatal("NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD is not set")
	}

	frontEndUrl := os.Getenv("FRONTEND_URL")
	if frontEndUrl == "" {
		log.Fatal("FRONTEND_URL is not set")
	}

	return &AppConfig{
		Port:                                 port,
		NotionAPIKey:                         apiKey,
		NotionKaimemoDatabaseInputID:         notionKaimemoDatabaseInputID,
		NotionKaimemoDatabaseSummaryRecordID: notionKaimemoDatabaseSummaryRecordID,
		AllowOrigins: []string{
			"http://localhost:5173", "http://localhost:4173", frontEndUrl,
		},
	}
}
