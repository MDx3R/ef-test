package gin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MDx3R/ef-test/internal/transport/http/response"
	"github.com/MDx3R/ef-test/internal/usecase"
	"github.com/MDx3R/ef-test/internal/usecase/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subService usecase.SubscriptionService
}

// Get godoc
// @Summary Получить подписку по ID
// @Description Возвращает подписку по заданному UUID
// @Tags subscriptions
// @Param id path string true "Subscription ID" Format(uuid)
// @Produce json
// @Success 200 {object} dto.SubscriptionResponse "Подписка найдена"
// @Failure 400 {object} response.ErrorResponse "Неверный UUID"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Get(ctx *gin.Context) {
	id, ok := h.parseUUIDParam(ctx, "id")
	if !ok {
		return
	}

	sub, err := h.subService.GetSubscription(id)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, sub)
}

// List godoc
// @Summary Список подписок
// @Description Возвращает список подписок с фильтрацией по параметрам
// @Tags subscriptions
// @Produce json
// @Param filter query dto.SubscriptionFilter false "Фильтры подписок"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации параметров запроса"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(ctx *gin.Context) {
	var filter dto.SubscriptionFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.handleValidationError(ctx, err)
		return
	}

	subs, err := h.subService.ListSubscriptions(filter)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, subs)
}

// Delete godoc
// @Summary Удалить подписку по ID
// @Description Удаляет подписку по UUID
// @Tags subscriptions
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 204 "Подписка удалена"
// @Failure 400 {object} response.ErrorResponse "Неверный UUID"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(ctx *gin.Context) {
	id, ok := h.parseUUIDParam(ctx, "id")
	if !ok {
		return
	}

	err := h.subService.DeleteSubscription(id)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// Create godoc
// @Summary Создать подписку
// @Description Создает новую подписку с данными из JSON
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequests true "Данные новой подписки"
// @Success 201 {object} response.IDResponse "ID созданной подписки"
// @Failure 400 {object} response.ErrorResponse "Неверный запрос"
// @Failure 422 {object} response.ValidationErrorResponse "Ошибка валидации"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(ctx *gin.Context) {
	var request dto.CreateSubscriptionRequests

	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		h.handleValidationError(ctx, err)
		return
	}

	id, err := h.subService.CreateSubscription(request)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.IDResponse{ID: id})
}

// Update godoc
// @Summary Обновить подписку
// @Description Обновляет подписку по UUID с данными из JSON
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Param subscription body dto.UpdateSubscriptionRequests true "Данные обновления подписки"
// @Success 204 "Подписка обновлена"
// @Failure 400 {object} response.ErrorResponse "Неверный UUID или данные запроса"
// @Failure 422 {object} response.ValidationErrorResponse "Ошибка валидации"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(ctx *gin.Context) {
	var request dto.UpdateSubscriptionRequests

	id, ok := h.parseUUIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := ctx.ShouldBindBodyWithJSON(&request); err != nil {
		h.handleValidationError(ctx, err)
		return
	}

	err := h.subService.UpdateSubscription(id, request)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// CalculateTotalCost godoc
// @Summary Рассчитать общую стоимость подписок
// @Description Возвращает общую стоимость подписок по фильтру
// @Tags subscriptions
// @Produce json
// @Param filter query dto.TotalCostFilter true "Фильтр для расчета стоимости"
// @Success 200 {object} response.IntResponse "Результат расчета стоимости"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации параметров запроса"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateTotalCost(ctx *gin.Context) {
	var request dto.TotalCostFilter

	if err := ctx.ShouldBindQuery(&request); err != nil {
		h.handleValidationError(ctx, err)
		return
	}

	result, err := h.subService.CalculateTotalCost(request)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.IntResponse{Value: result})
}

func (h *SubscriptionHandler) parseUUIDParam(ctx *gin.Context, param string) (uuid.UUID, bool) {
	idStr := ctx.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(ctx, http.StatusBadRequest, fmt.Errorf("uuid not valid: %s", ctx.Param("id")))
		return uuid.Nil, false
	}
	return id, true
}

func (h *SubscriptionHandler) handleServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		h.respondError(ctx, http.StatusNotFound, err)
	default:
		h.respondError(ctx, http.StatusInternalServerError, err)
	}
}

func (h *SubscriptionHandler) handleValidationError(ctx *gin.Context, err error) {
	var verr validator.ValidationErrors

	if errors.As(err, &verr) {
		errorsMap := h.buildMap(verr)
		h.respondValidationError(ctx, errorsMap)
		return
	}

	h.respondError(ctx, http.StatusBadRequest, err)
}

func (h *SubscriptionHandler) buildMap(verr validator.ValidationErrors) map[string]string {
	errorsMap := make(map[string]string)
	for _, fe := range verr {
		errorsMap[fe.Field()] = fmt.Sprintf("field '%s' validation failed on '%s' tag", fe.Field(), fe.Tag())
	}
	return errorsMap
}

func (h *SubscriptionHandler) respondError(ctx *gin.Context, code int, err error) {
	ctx.AbortWithStatusJSON(code, response.ErrorResponse{Error: err.Error()})
}

func (h *SubscriptionHandler) respondValidationError(ctx *gin.Context, errMap map[string]string) {
	ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response.ValidationErrorResponse{
		Error:  "validation error",
		Fields: errMap,
	})
}

func NewSubscriptionHandler(subService usecase.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subService: subService,
	}
}
