package service

import (
	"errors"
	"template-echo-notion-integration/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKaimemoRepository struct {
	mock.Mock
}

func (m *MockKaimemoRepository) FetchKaimemo() ([]model.KaimemoResponse, error) {
	args := m.Called()
	return args.Get(0).([]model.KaimemoResponse), args.Error(1)
}

func (m *MockKaimemoRepository) InsertKaimemo(req model.CreateKaimemoRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockKaimemoRepository) RemoveKaimemo(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestKaimemoService_FetchKaimemo(t *testing.T) {
	mockRepo := new(MockKaimemoRepository)
	service := NewKaimemoService(mockRepo)

	t.Run("success fetch kaimemo", func(t *testing.T) {
		expected := []model.KaimemoResponse{
			{ID: "1", Name: "Test 1", Tag: "Content 1", Done: false},
			{ID: "2", Name: "Test 2", Tag: "Content 2", Done: false},
		}
		mockRepo.On("FetchKaimemo").Return(expected, nil).Once()

		result, err := service.FetchKaimemo()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error fetch kaimemo", func(t *testing.T) {
		mockRepo.On("FetchKaimemo").Return([]model.KaimemoResponse{}, errors.New("fetch error")).Once()

		result, err := service.FetchKaimemo()

		assert.Error(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestKaimemoService_CreateKaimemo(t *testing.T) {
	mockRepo := new(MockKaimemoRepository)
	service := NewKaimemoService(mockRepo)

	t.Run("success create kaimemo", func(t *testing.T) {
		req := model.CreateKaimemoRequest{
			Name: "Test Title",
			Tag:  "Test Content",
		}
		mockRepo.On("InsertKaimemo", req).Return(nil).Once()

		err := service.CreateKaimemo(req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error create kaimemo", func(t *testing.T) {
		req := model.CreateKaimemoRequest{
			Name: "Test Title",
			Tag:  "Test Content",
		}
		mockRepo.On("InsertKaimemo", req).Return(errors.New("insert error")).Once()

		err := service.CreateKaimemo(req)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestKaimemoService_RemoveKaimemo(t *testing.T) {
	mockRepo := new(MockKaimemoRepository)
	service := NewKaimemoService(mockRepo)

	t.Run("success remove kaimemo", func(t *testing.T) {
		id := "test-id"
		mockRepo.On("RemoveKaimemo", id).Return(nil).Once()

		err := service.RemoveKaimemo(id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error remove kaimemo", func(t *testing.T) {
		id := "test-id"
		mockRepo.On("RemoveKaimemo", id).Return(errors.New("remove error")).Once()

		err := service.RemoveKaimemo(id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
