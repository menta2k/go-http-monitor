package result

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleLatest(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	cr, err := h.service.Latest(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		response.WriteError(w, http.StatusNotFound, "no check results yet")
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get latest result")
		return
	}
	response.WriteJSON(w, http.StatusOK, cr)
}

func (h *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	results, err := h.service.History(r.Context(), id, limit, offset)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get history")
		return
	}
	if results == nil {
		results = []domain.CheckResult{}
	}
	response.WriteJSON(w, http.StatusOK, results)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
