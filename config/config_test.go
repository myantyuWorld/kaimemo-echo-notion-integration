package config

import (
	"os"
	"testing"
)

func TestLoadConfigWithCustomPort(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("NOTION_API_KEY", "test-api-key")
	os.Setenv("NOTION_DATABASE_ID", "test-database-id")
	os.Setenv("FRONTEND_URL", "http://test-frontend.com")

	config := LoadConfig()

	if config.Port != "8080" {
		t.Errorf("expected Port to be 8080, got %s", config.Port)
	}

	os.Unsetenv("PORT")
	os.Unsetenv("NOTION_API_KEY")
	os.Unsetenv("NOTION_DATABASE_ID")
	os.Unsetenv("FRONTEND_URL")
}

func TestLoadConfigWithMultipleFrontendURLs(t *testing.T) {
	os.Setenv("NOTION_API_KEY", "test-api-key")
	os.Setenv("NOTION_DATABASE_ID", "test-database-id")
	os.Setenv("FRONTEND_URL", "http://test1.com,http://test2.com")

	config := LoadConfig()

	expectedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:4173",
		"http://test1.com",
		"http://test2.com",
	}

	if len(config.AllowOrigins) != len(expectedOrigins) {
		t.Errorf("expected %d origins, got %d", len(expectedOrigins), len(config.AllowOrigins))
	}

	for i, origin := range config.AllowOrigins {
		if origin != expectedOrigins[i] {
			t.Errorf("expected AllowOrigins[%d] to be %s, got %s", i, expectedOrigins[i], origin)
		}
	}

	os.Unsetenv("NOTION_API_KEY")
	os.Unsetenv("NOTION_DATABASE_ID")
	os.Unsetenv("FRONTEND_URL")
}

func TestLoadConfigWithEmptyPort(t *testing.T) {
	os.Setenv("PORT", "")
	os.Setenv("NOTION_API_KEY", "test-api-key")
	os.Setenv("NOTION_DATABASE_ID", "test-database-id")
	os.Setenv("FRONTEND_URL", "http://test-frontend.com")

	config := LoadConfig()

	if config.Port != "3000" {
		t.Errorf("expected Port to be default value 3000 when empty, got %s", config.Port)
	}

	os.Unsetenv("PORT")
	os.Unsetenv("NOTION_API_KEY")
	os.Unsetenv("NOTION_DATABASE_ID")
	os.Unsetenv("FRONTEND_URL")
}
