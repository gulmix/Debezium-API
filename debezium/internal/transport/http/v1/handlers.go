package v1

import (
	"context"
	"debezium_server/internal/models"
	"encoding/json"
	"net/http"
)

type UserService interface {
	GetUsers(ctx context.Context, offset, limit int) ([]models.User, error)
}

type HandlerFacade struct {
	service UserService
}

func NewHandlerFacade(service UserService) *HandlerFacade {
	return &HandlerFacade{service: service}
}

func (h *HandlerFacade) GetUsers(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10

	users, err := h.service.GetUsers(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
