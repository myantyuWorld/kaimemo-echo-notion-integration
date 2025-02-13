package repository

import (
	"testing"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
)

func TestNewNotionRepository(t *testing.T) {
	testCases := []struct {
		name       string
		apiKey     string
		databaseID string
	}{
		{
			name:       "should create new repository with valid credentials",
			apiKey:     "test-api-key",
			databaseID: "test-database-id",
		},
		{
			name:       "should create new repository with empty credentials",
			apiKey:     "",
			databaseID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewNotionRepository(tc.apiKey, tc.databaseID)

			assert.NotNil(t, repo)

			kaimemoRepo, ok := repo.(*kaimemoRepository)
			assert.True(t, ok)
			assert.NotNil(t, kaimemoRepo.client)
			assert.Equal(t, tc.databaseID, kaimemoRepo.databaseID)
			assert.Equal(t, notionapi.Token(tc.apiKey), kaimemoRepo.client.Token)
		})
	}
}
