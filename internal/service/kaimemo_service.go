package service

import (
	"template-echo-notion-integration/internal/model"
	"template-echo-notion-integration/internal/repository"
)

type kaimemoService struct {
	repo repository.KaimemoRepository
}

// CreateKaimemo implements KaimemoService.
func (k *kaimemoService) CreateKaimemo(req model.CreateKaimemoRequest) error {
	return k.repo.InsertKaimemo(req)
}

// FetchKaimemo implements KaimemoService.
func (k *kaimemoService) FetchKaimemo() ([]model.KaimemoResponse, error) {
	return k.repo.FetchKaimemo()
}

// RemoveKaimemo implements KaimemoService.
func (k *kaimemoService) RemoveKaimemo(id string) error {
	return k.repo.RemoveKaimemo(id)
}

type KaimemoService interface {
	FetchKaimemo() ([]model.KaimemoResponse, error)
	CreateKaimemo(req model.CreateKaimemoRequest) error
	RemoveKaimemo(id string) error
}

func NewKaimemoService(repo repository.KaimemoRepository) KaimemoService {
	return &kaimemoService{repo: repo}
}
