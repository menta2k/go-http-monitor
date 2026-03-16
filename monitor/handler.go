package monitor

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/response"
)

type Handler struct {
	service   *Service
	onChanged func()
}

func NewHandler(service *Service, onChanged func()) *Handler {
	return &Handler{service: service, onChanged: onChanged}
}

type createRequest struct {
	URL             string `json:"url"`
	ExpectedStatus  int    `json:"expected_status"`
	BodyContains    string `json:"body_contains"`
	IntervalSeconds int    `json:"interval_seconds"`
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	monitors, err := h.service.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to list monitors")
		return
	}
	if monitors == nil {
		monitors = []domain.Monitor{}
	}
	response.WriteJSON(w, http.StatusOK, monitors)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	m, err := h.service.Get(r.Context(), id)
	if errors.Is(err, domain.ErrMonitorNotFound) {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get monitor")
		return
	}
	response.WriteJSON(w, http.StatusOK, m)
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	m, err := h.service.Create(r.Context(), req.URL, req.ExpectedStatus, req.BodyContains, req.IntervalSeconds)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) || errors.Is(err, domain.ErrInvalidStatusCode) || errors.Is(err, domain.ErrInvalidInterval) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to create monitor")
		return
	}

	if h.onChanged != nil {
		h.onChanged()
	}
	response.WriteJSON(w, http.StatusCreated, m)
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	m, err := h.service.Update(r.Context(), id, req.URL, req.ExpectedStatus, req.BodyContains, req.IntervalSeconds)
	if err != nil {
		if errors.Is(err, domain.ErrMonitorNotFound) {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidURL) || errors.Is(err, domain.ErrInvalidStatusCode) || errors.Is(err, domain.ErrInvalidInterval) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to update monitor")
		return
	}

	if h.onChanged != nil {
		h.onChanged()
	}
	response.WriteJSON(w, http.StatusOK, m)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrMonitorNotFound) {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to delete monitor")
		return
	}

	if h.onChanged != nil {
		h.onChanged()
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
