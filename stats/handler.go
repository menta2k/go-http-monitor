package stats

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sko/go-http-monitor/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	monitorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	period := ParsePeriod(r.URL.Query().Get("period"))

	summary, err := h.service.GetSummary(r.Context(), monitorID, period)
	if err != nil {
		log.Printf("[stats] summary error for monitor %d: %v", monitorID, err)
		response.WriteError(w, http.StatusInternalServerError, "failed to get summary")
		return
	}
	response.WriteJSON(w, http.StatusOK, summary)
}

func (h *Handler) HandleTimeline(w http.ResponseWriter, r *http.Request) {
	monitorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	period := ParsePeriod(r.URL.Query().Get("period"))
	buckets, _ := strconv.Atoi(r.URL.Query().Get("buckets"))
	buckets = ParseBuckets(buckets)

	points, err := h.service.GetTimeline(r.Context(), monitorID, period, buckets)
	if err != nil {
		log.Printf("[stats] timeline error for monitor %d: %v", monitorID, err)
		response.WriteError(w, http.StatusInternalServerError, "failed to get timeline")
		return
	}
	response.WriteJSON(w, http.StatusOK, points)
}

func (h *Handler) HandleStatusCodes(w http.ResponseWriter, r *http.Request) {
	monitorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	period := ParsePeriod(r.URL.Query().Get("period"))

	codes, err := h.service.GetStatusCodes(r.Context(), monitorID, period)
	if err != nil {
		log.Printf("[stats] status codes error for monitor %d: %v", monitorID, err)
		response.WriteError(w, http.StatusInternalServerError, "failed to get status codes")
		return
	}
	response.WriteJSON(w, http.StatusOK, codes)
}

func (h *Handler) HandleStatusCodeTimeline(w http.ResponseWriter, r *http.Request) {
	monitorID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid monitor id")
		return
	}

	period := ParsePeriod(r.URL.Query().Get("period"))
	buckets, _ := strconv.Atoi(r.URL.Query().Get("buckets"))
	buckets = ParseBuckets(buckets)

	points, err := h.service.GetStatusCodeTimeline(r.Context(), monitorID, period, buckets)
	if err != nil {
		log.Printf("[stats] status code timeline error for monitor %d: %v", monitorID, err)
		response.WriteError(w, http.StatusInternalServerError, "failed to get status code timeline")
		return
	}
	response.WriteJSON(w, http.StatusOK, points)
}
