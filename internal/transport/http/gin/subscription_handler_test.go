package gin_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	logruslogger "github.com/MDx3R/ef-test/internal/infra/logger"
	handlers "github.com/MDx3R/ef-test/internal/transport/http/gin"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	mock_usecase "github.com/MDx3R/ef-test/internal/usecase/mocks"
	"github.com/MDx3R/ef-test/internal/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger = logruslogger.NewLogger()

func setupRouterAndHandler(t *testing.T) (*gin.Engine, *mock_usecase.MockSubscriptionService) {
	gin.SetMode(gin.TestMode)

	mockService := mock_usecase.NewMockSubscriptionService(t)
	handler := handlers.NewSubscriptionHandler(mockService, logger)

	r := gin.New()
	r.GET("", handler.List)
	r.GET("/:id", handler.Get)
	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
	r.GET("/total", handler.CalculateTotalCost)

	return r, mockService
}

func makeTestSubscriptionResponse(t *testing.T) dto.SubscriptionResponse {
	id := uuid.New()
	return dto.SubscriptionResponse{
		ID:          id,
		ServiceName: "test_service",
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)),
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

	startDate := model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))
	request := dto.CreateSubscriptionRequests{
		ServiceName: "test_service",
		Price:       100,
		UserID:      uuid.New(),
		StartDate:   &startDate,
		EndDate:     nil,
	}

	jsonBody := fmt.Sprintf(
		`{"service_name":"%v", "price":%v, "user_id":"%v", "start_date":"08-2025"}`,
		request.ServiceName,
		request.Price,
		request.UserID,
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

func TestSubscriptionHandler_Create_InvalidJSON_Validation(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	tests := []struct {
		name       string
		jsonBody   string
		expectCode int
		expectErr  string
	}{
		{
			name:       "missing service_name",
			jsonBody:   `{"price":100,"user_id":"` + uuid.New().String() + `","start_date":"08-2025"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "ServiceName",
		},
		{
			name:       "missing price",
			jsonBody:   `{"service_name":"Test","user_id":"` + uuid.New().String() + `","start_date":"08-2025"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "Price",
		},
		{
			name:       "missing user_id",
			jsonBody:   `{"service_name":"Test","price":100,"start_date":"08-2025"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "UserID",
		},
		{
			name:       "missing start_date",
			jsonBody:   `{"service_name":"Test","price":100,"user_id":"` + uuid.New().String() + `"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "StartDate",
		},
		{
			name:       "invalid month-year format",
			jsonBody:   `{"service_name":"Test","price":100,"user_id":"` + uuid.New().String() + `","start_date":"2025-08-01"}`,
			expectCode: http.StatusBadRequest,
			expectErr:  "error",
		},
		{
			name:       "invalid json syntax",
			jsonBody:   `{invalid_json}`,
			expectCode: http.StatusBadRequest,
			expectErr:  "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectErr)
		})
	}
}

func TestSubscriptionHandler_Update_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	id := uuid.New()
	startDate := model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))
	request := dto.UpdateSubscriptionRequests{
		ServiceName: "test_service",
		Price:       100,
		StartDate:   &startDate,
		EndDate:     nil,
	}

	jsonBody := fmt.Sprintf(
		`{"service_name":"%v", "price":%v, "start_date":"08-2025"}`,
		request.ServiceName,
		request.Price,
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

func TestSubscriptionHandler_Update_InvalidJSON_Validation(t *testing.T) {
	router, _ := setupRouterAndHandler(t)

	tests := []struct {
		name       string
		jsonBody   string
		expectCode int
		expectErr  string
	}{
		{
			name:       "missing service_name",
			jsonBody:   `{"price":100,"start_date":"08-2025"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "ServiceName",
		},
		{
			name:       "missing price",
			jsonBody:   `{"service_name":"Test","start_date":"08-2025"}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "Price",
		},
		{
			name:       "missing start_date",
			jsonBody:   `{"service_name":"Test","price":100}`,
			expectCode: http.StatusUnprocessableEntity,
			expectErr:  "StartDate",
		},
		{
			name:       "invalid month-year format",
			jsonBody:   `{"service_name":"Test","price":100,"start_date":"2025-08-01"}`,
			expectCode: http.StatusBadRequest,
			expectErr:  "error",
		},
		{
			name:       "invalid json syntax",
			jsonBody:   `{invalid_json}`,
			expectCode: http.StatusBadRequest,
			expectErr:  "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			req := httptest.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(tc.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectErr)
		})
	}
}

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

func TestSubscriptionHandler_Get_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	subID := uuid.New()
	mockService.On("GetSubscription", subID).Return(dto.SubscriptionResponse{}, errors.New("service failure"))

	req := httptest.NewRequest(http.MethodGet, "/"+subID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_List_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	mockService.On("ListSubscriptions", mock.Anything).Return(nil, errors.New("service failure"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_Create_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	jsonBody := fmt.Sprintf(
		`{"service_name":"%v", "price":%v, "user_id":"%v", "start_date":"08-2025"}`,
		"test_service",
		100,
		uuid.New().String(),
	)

	mockService.On("CreateSubscription", mock.Anything).Return(uuid.Nil, errors.New("service failure"))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_Update_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	subID := uuid.New()
	jsonBody := `{"service_name":"Updated","price":200,"start_date":"09-2025"}`
	mockService.On("UpdateSubscription", subID, mock.Anything).Return(errors.New("service failure"))

	req := httptest.NewRequest(http.MethodPut, "/"+subID.String(), strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_Delete_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	subID := uuid.New()
	mockService.On("DeleteSubscription", subID).Return(errors.New("service failure"))

	req := httptest.NewRequest(http.MethodDelete, "/"+subID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_TotalCost_Success(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	userID := uuid.New()
	periodStart := model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))
	periodEnd := model.NewMonthYear(time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC))
	request := dto.TotalCostFilter{
		UserID:      userID.String(),
		ServiceName: "test_service",
		PeriodStart: &periodStart,
		PeriodEnd:   &periodEnd,
	}

	query := fmt.Sprintf(
		`/total?user_id=%s&service_name=%s&period_start="08-2025"&period_end="10-2025"`,
		userID.String(),
		"test_service",
	)

	mockService.On(
		"CalculateTotalCost",
		request,
	).Return(150, nil)

	req := httptest.NewRequest(http.MethodGet, query, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "150")
	mockService.AssertExpectations(t)
}

func TestSubscriptionHandler_TotalCost_ServiceError(t *testing.T) {
	router, mockService := setupRouterAndHandler(t)

	userID := uuid.New()
	periodStart := model.NewMonthYear(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC))
	periodEnd := model.NewMonthYear(time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC))
	request := dto.TotalCostFilter{
		UserID:      userID.String(),
		ServiceName: "test_service",
		PeriodStart: &periodStart,
		PeriodEnd:   &periodEnd,
	}

	query := fmt.Sprintf(
		`/total?user_id=%s&service_name=%s&period_start="08-2025"&period_end="10-2025"`,
		userID.String(),
		"test_service",
	)

	mockService.On(
		"CalculateTotalCost",
		request,
	).Return(0, errors.New("service failure"))

	req := httptest.NewRequest(http.MethodGet, query, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
