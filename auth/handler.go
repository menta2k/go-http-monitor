package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sko/go-http-monitor/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		response.WriteError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	token, err := h.service.Authenticate(req.Username, req.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "authentication failed")
		return
	}

	response.WriteJSON(w, http.StatusOK, loginResponse{Token: token})
}
