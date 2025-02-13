package repository

import (
	"context"
	"log"
	"template-echo-notion-integration/internal/model"

	"github.com/jomei/notionapi"
)

type kaimemoRepository struct {
	client     *notionapi.Client
	query      *notionapi.DatabaseQueryRequest
	databaseID string
}

// FetchKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) FetchKaimemo() ([]model.KaimemoResponse, error) {
	resp, err := k.client.Database.Query(context.Background(), notionapi.DatabaseID(k.databaseID), k.query)
	if err != nil {
		log.Printf("failed to notion query database: %v", err)
		return nil, err
	}

	var kaimemoResponses []model.KaimemoResponse
	for _, result := range resp.Results {
		properties := result.Properties

		data := model.KaimemoResponse{}
		data.ID = string(result.ID)
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

	return kaimemoResponses, nil
}

// InsertKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) InsertKaimemo(req model.CreateKaimemoRequest) error {
	_, err := k.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(k.databaseID), // 既存のデータベースID
		},
		Properties: notionapi.Properties{
			"name": &notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: req.Name,
						},
					},
				},
			},
			"tag": &notionapi.SelectProperty{
				Select: notionapi.Option{
					Name: req.Tag,
				},
			},
		},
	})

	if err != nil {
		log.Printf("failed to notion create page: %v", err)
		return err
	}

	return nil
}

// RemoveKaimemo implements KaimemoRepository.
func (k *kaimemoRepository) RemoveKaimemo(id string) error {
	_, err := k.client.Page.Update(context.Background(), notionapi.PageID(id), &notionapi.PageUpdateRequest{
		Archived: true,
	})

	if err != nil {
		log.Printf("failed to notion update page: %v", err)
		return err
	}

	return nil
}

type KaimemoRepository interface {
	FetchKaimemo() ([]model.KaimemoResponse, error)
	InsertKaimemo(req model.CreateKaimemoRequest) error
	RemoveKaimemo(id string) error
}

func NewNotionRepository(apiKey, databaseID string) KaimemoRepository {
	client := notionapi.NewClient(notionapi.Token(apiKey))
	query := &notionapi.DatabaseQueryRequest{}

	return &kaimemoRepository{client: client, databaseID: databaseID, query: query}
}
