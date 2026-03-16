package notification

import (
	"encoding/json"
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

type createRequest struct {
	Type    domain.NotificationType `json:"type"`
	Target  string                  `json:"target"`
	Enabled bool                    `json:"enabled"`
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	monitorID, err := parseMonitorID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	notifications, err := h.service.ListByMonitor(r.Context(), monitorID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}
	if notifications == nil {
		notifications = []domain.Notification{}
	}
	response.WriteJSON(w, http.StatusOK, notifications)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id, err := parseNotifID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	n, err := h.service.Get(r.Context(), id)
	if errors.Is(err, domain.ErrNotificationNotFound) {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get notification")
		return
	}
	response.WriteJSON(w, http.StatusOK, n)
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	monitorID, err := parseMonitorID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := h.service.Create(r.Context(), monitorID, req.Type, req.Target, req.Enabled)
	if err != nil {
		if isValidationError(err) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to create notification")
		return
	}
	response.WriteJSON(w, http.StatusCreated, n)
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseNotifID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	n, err := h.service.Update(r.Context(), id, req.Type, req.Target, req.Enabled)
	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		if isValidationError(err) {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to update notification")
		return
	}
	response.WriteJSON(w, http.StatusOK, n)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseNotifID(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			response.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to delete notification")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseMonitorID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func parseNotifID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("nid"), 10, 64)
}

func isValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidNotificationType) ||
		errors.Is(err, domain.ErrInvalidNotificationTarget) ||
		errors.Is(err, domain.ErrMonitorNotFound)
}
