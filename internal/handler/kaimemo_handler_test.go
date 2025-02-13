package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"template-echo-notion-integration/internal/model"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKaimemoService struct {
	mock.Mock
}

func (m *MockKaimemoService) FetchKaimemo() ([]model.KaimemoResponse, error) {
	args := m.Called()
	return args.Get(0).([]model.KaimemoResponse), args.Error(1)
}

func (m *MockKaimemoService) CreateKaimemo(req model.CreateKaimemoRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockKaimemoService) RemoveKaimemo(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestKaimemoHandler_CreateKaimemo(t *testing.T) {
	mockService := new(MockKaimemoService)
	h := NewKaimemoHandler(mockService)
	e := echo.New()

	t.Run("invalid request body", func(t *testing.T) {
		reqBody := `{"invalid json`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("empty request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("CreateKaimemo", model.CreateKaimemoRequest{}).Return(nil).Once()
		err := h.CreateKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestKaimemoHandler_FetchKaimemo(t *testing.T) {
	mockService := new(MockKaimemoService)
	h := NewKaimemoHandler(mockService)
	e := echo.New()

	t.Run("service returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("FetchKaimemo").Return([]model.KaimemoResponse{}, nil).Once()
		err := h.FetchKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]", strings.TrimSpace(rec.Body.String()))
	})

	t.Run("service returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("FetchKaimemo").Return([]model.KaimemoResponse{}, errors.New("service error")).Once()
		err := h.FetchKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestKaimemoHandler_RemoveKaimemo(t *testing.T) {
	mockService := new(MockKaimemoService)
	h := NewKaimemoHandler(mockService)
	e := echo.New()

	t.Run("empty id parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.RemoveKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("test-id")

		mockService.On("RemoveKaimemo", "test-id").Return(errors.New("service error")).Once()
		err := h.RemoveKaimemo(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
