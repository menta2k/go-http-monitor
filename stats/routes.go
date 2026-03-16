package stats

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/monitors/{id}/stats", h.HandleSummary)
	mux.HandleFunc("GET /api/monitors/{id}/timeline", h.HandleTimeline)
}
