package service

import (
	"template-echo-notion-integration/internal/model"
	"template-echo-notion-integration/internal/repository"
)

type kaimemoService struct {
	repo repository.KaimemoRepository
}

// CreateKaimemoAmount implements KaimemoService.
func (k *kaimemoService) CreateKaimemoAmount(req model.CreateKaimemoAmountRequest) error {
	return k.repo.InsertKaimemoAmount(req)
}

// FetchKaimemoSummaryRecord implements KaimemoService.
func (k *kaimemoService) FetchKaimemoSummaryRecord() (model.KaimemoSummaryResponse, error) {
	res, err := k.repo.FetchKaimemoAmountRecords()
	if err != nil {
		return model.KaimemoSummaryResponse{}, err
	}

	return model.KaimemoSummaryResponse{
		MonthlySummaries: res.GroupByMonth(),
		WeeklySummaries:  res.GroupByWeek(),
	}, nil
}

// RemoveKaimemoAmount implements KaimemoService.
func (k *kaimemoService) RemoveKaimemoAmount(id string) error {
	return k.repo.RemoveKaimemoAmount(id)
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
	FetchKaimemoSummaryRecord() (model.KaimemoSummaryResponse, error)
	CreateKaimemoAmount(req model.CreateKaimemoAmountRequest) error
	RemoveKaimemoAmount(id string) error
}

func NewKaimemoService(repo repository.KaimemoRepository) KaimemoService {
	return &kaimemoService{repo: repo}
}
