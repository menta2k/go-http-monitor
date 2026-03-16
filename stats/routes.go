package stats

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/monitors/{id}/stats", h.HandleSummary)
	mux.HandleFunc("GET /api/monitors/{id}/timeline", h.HandleTimeline)
	mux.HandleFunc("GET /api/monitors/{id}/status-codes", h.HandleStatusCodes)
	mux.HandleFunc("GET /api/monitors/{id}/status-code-timeline", h.HandleStatusCodeTimeline)
}
