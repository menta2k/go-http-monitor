package result

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/monitors/{id}/status", h.HandleLatest)
	mux.HandleFunc("GET /api/monitors/{id}/history", h.HandleHistory)
}
