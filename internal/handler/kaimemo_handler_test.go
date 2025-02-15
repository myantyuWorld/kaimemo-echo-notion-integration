package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"template-echo-notion-integration/internal/model"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKaimemoService struct {
	mock.Mock
}

// CreateKaimemoAmount implements service.KaimemoService.
func (m *MockKaimemoService) CreateKaimemoAmount(req model.CreateKaimemoAmountRequest) error {
	panic("unimplemented")
}

// FetchKaimemoSummaryRecord implements service.KaimemoService.
func (m *MockKaimemoService) FetchKaimemoSummaryRecord() ([]model.WeeklySummary, error) {
	args := m.Called()
	return args.Get(0).([]model.WeeklySummary), args.Error(1)
}

// RemoveKaimemoAmount implements service.KaimemoService.
func (m *MockKaimemoService) RemoveKaimemoAmount(id string) error {
	panic("unimplemented")
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
func TestKaimemoHandler_FetchKaimemoSummaryRecord(t *testing.T) {
	mockService := new(MockKaimemoService)
	h := NewKaimemoHandler(mockService)
	e := echo.New()

	t.Run("successful fetch with data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedSummary := []model.WeeklySummary{
			{
				WeekStart:   "2023-05-28",
				WeekEnd:     "2023-06-03",
				TotalAmount: 3000,
				Items: []model.KaimemoAmount{
					{ID: "", Date: "2023-05-30", Tag: "食費", Amount: 1000},
					{ID: "", Date: "2023-06-01", Tag: "食費", Amount: 2000},
				},
			},
			{
				WeekStart:   "2023-06-04",
				WeekEnd:     "2023-06-10",
				TotalAmount: 3000,
				Items: []model.KaimemoAmount{
					{ID: "", Date: "2023-06-05", Tag: "日用品", Amount: 3000},
				},
			},
		}

		expectedJson := `[
				{
					"weekStart": "2023-05-28",
					"weekEnd": "2023-06-03",
					"totalAmount": 3000,
					"items": [
						{
							"id" : "",
							"tag": "食費",
							"date": "2023-05-30",
							"amount": 1000
						},
						{
							"id" : "",
							"tag": "食費",
							"date": "2023-06-01",
							"amount": 2000
						}
					]
				},
				{
					"weekStart": "2023-06-04",
					"weekEnd": "2023-06-10",
					"totalAmount": 3000,
					"items": [
						{
							"id" : "",
							"tag": "日用品",
							"date": "2023-06-05",
							"amount": 3000
						}
					]
				}
			]`

		mockService.On("FetchKaimemoSummaryRecord").Return(expectedSummary, nil).Once()
		err := h.FetchKaimemoSummaryRecord(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		spew.Dump(rec.Body.String())
		assert.JSONEq(t, expectedJson, rec.Body.String())
	})
	t.Run("successful fetch with empty data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("FetchKaimemoSummaryRecord").Return([]model.WeeklySummary{}, nil).Once()
		err := h.FetchKaimemoSummaryRecord(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]", strings.TrimSpace(rec.Body.String()))
	})

	t.Run("service returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("FetchKaimemoSummaryRecord").Return([]model.WeeklySummary{}, errors.New("database error")).Once()
		err := h.FetchKaimemoSummaryRecord(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to fetch kaimemo summary record")
	})
}
