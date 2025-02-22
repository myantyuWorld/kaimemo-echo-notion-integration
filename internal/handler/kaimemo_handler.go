package handler

import (
	"net/http"
	"template-echo-notion-integration/internal/model"
	"template-echo-notion-integration/internal/service"

	"github.com/labstack/echo/v4"
)

type kaimemoHandler struct {
	service service.KaimemoService
}

// CreateKaimemoAmount implements KaimemoHandler.
func (k *kaimemoHandler) CreateKaimemoAmount(c echo.Context) error {
	req := model.CreateKaimemoAmountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := k.service.CreateKaimemoAmount(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create kaimemo amount",
		})
	}

	return c.NoContent(http.StatusCreated)
}

// FetchKaimemoSummaryRecord implements KaimemoHandler.
func (k *kaimemoHandler) FetchKaimemoSummaryRecord(c echo.Context) error {
	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "TempUserID is required",
		})
	}

	res, err := k.service.FetchKaimemoSummaryRecord(tempUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch kaimemo summary record",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// RemoveKaimemoAmount implements KaimemoHandler.
func (k *kaimemoHandler) RemoveKaimemoAmount(c echo.Context) error {
	req := model.RemoveKaimemoAmountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	if err := k.service.RemoveKaimemoAmount(id, req.TempUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove kaimemo",
		})
	}

	return c.NoContent(http.StatusOK)
}

// CreateKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) CreateKaimemo(c echo.Context) error {
	req := model.CreateKaimemoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	if err := k.service.CreateKaimemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create kaimemo",
		})
	}

	return c.NoContent(http.StatusCreated)
}

// FetchKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) FetchKaimemo(c echo.Context) error {
	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "TempUserID is required",
		})
	}

	res, err := k.service.FetchKaimemo(tempUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch kaimemo",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// RemoveKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) RemoveKaimemo(c echo.Context) error {
	req := model.RemoveKaimemoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	if err := k.service.RemoveKaimemo(id, req.TempUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove kaimemo",
		})
	}

	return c.NoContent(http.StatusOK)
}

type KaimemoHandler interface {
	FetchKaimemo(c echo.Context) error
	CreateKaimemo(c echo.Context) error
	RemoveKaimemo(c echo.Context) error
	FetchKaimemoSummaryRecord(c echo.Context) error
	CreateKaimemoAmount(c echo.Context) error
	RemoveKaimemoAmount(c echo.Context) error
}

func NewKaimemoHandler(service service.KaimemoService) KaimemoHandler {
	return &kaimemoHandler{service: service}
}
