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
