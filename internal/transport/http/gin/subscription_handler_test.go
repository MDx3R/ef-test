package gin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	handlers "github.com/MDx3R/ef-test/internal/transport/http/gin"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	mock_usecase "github.com/MDx3R/ef-test/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRouterAndHandler(t *testing.T) (*gin.Engine, *mock_usecase.MockSubscriptionService) {
	gin.SetMode(gin.TestMode)

	mockService := mock_usecase.NewMockSubscriptionService(t)
	handler := handlers.NewSubscriptionHandler(mockService)

	r := gin.New()
	r.GET("", handler.List)
	r.GET("/:id", handler.Get)
	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
	r.GET("/total", handler.Delete)

	return r, mockService
}

func makeTestSubscriptionResponse(t *testing.T) dto.SubscriptionResponse {
	id := uuid.New()
	return dto.SubscriptionResponse{
		ID:          id,
		ServiceName: "test_service",
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}
}

func TestSubscriptionHandler_Get_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	sub := makeTestSubscriptionResponse(t)
	id := sub.ID

	mockService.On("GetSubscription", id).Return(sub, nil)

	req := httptest.NewRequest(http.MethodGet, "/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), fmt.Sprintf(`"id":"%s"`, id.String()))
	assert.Contains(t, w.Body.String(), `"service_name":"test_service"`)
	assert.Contains(t, w.Body.String(), `"price":100`)
	assert.Contains(t, w.Body.String(), fmt.Sprintf(`"user_id":"%s"`, sub.UserID.String()))
}

func TestSubscriptionHandler_Get_InvalidUUID(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/not-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "uuid not valid")
}

func TestSubscriptionHandler_Get_NotFound(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	sub := makeTestSubscriptionResponse(t)
	id := sub.ID

	mockService.On("GetSubscription", id).Return(dto.SubscriptionResponse{}, usecase.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_List_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	subs := []dto.SubscriptionResponse{
		makeTestSubscriptionResponse(t),
		makeTestSubscriptionResponse(t),
	}

	filter := dto.SubscriptionFilter{Page: 1, PageSize: 20}

	mockService.On("ListSubscriptions", filter).Return(subs, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), subs[0].ID.String())
	assert.Contains(t, w.Body.String(), subs[1].ID.String())
}

// TODO: Filter DTO tests

func TestSubscriptionHandler_Create_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	request := dto.CreateSubscriptionRequests{
		ServiceName: "test_service",
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}

	jsonBody := fmt.Sprintf(
		`{"service_name":"%v", "price":%v, "user_id":"%v", "start_date":"%v"}`,
		request.ServiceName,
		request.Price,
		request.UserID,
		request.StartDate.Format(time.RFC3339),
	)

	newID := uuid.New()
	mockService.On("CreateSubscription", request).Return(newID, nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
	assert.Contains(t, w.Body.String(), newID.String())
}

func TestSubscriptionHandler_Create_InvalidJSON(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	jsonBody := `{invalid_json}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TODO: Create Validation tests

func TestSubscriptionHandler_Update_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	id := uuid.New()
	request := dto.UpdateSubscriptionRequests{
		ServiceName: "test_service",
		Price:       100,
		StartDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     nil,
	}

	jsonBody := fmt.Sprintf(
		`{"service_name":"%v", "price":%v, "start_date":"%v"}`,
		request.ServiceName,
		request.Price,
		request.StartDate.Format(time.RFC3339),
	)

	mockService.On("UpdateSubscription", id, request).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_Update_InvalidUUID(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/not-uuid", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestSubscriptionHandler_Update_InvalidJSON(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	id := uuid.New()
	jsonBody := `{invalid_json}`

	req := httptest.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TODO: Update Validation tests

func TestSubscriptionHandler_Delete_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	id := uuid.New()
	mockService.On("DeleteSubscription", id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_Delete_InvalidUUID(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/not-uuid", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TODO: Service Error tests
