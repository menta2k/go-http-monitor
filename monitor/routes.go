package monitor

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/monitors", h.HandleList)
	mux.HandleFunc("POST /api/monitors", h.HandleCreate)
	mux.HandleFunc("GET /api/monitors/{id}", h.HandleGet)
	mux.HandleFunc("PUT /api/monitors/{id}", h.HandleUpdate)
	mux.HandleFunc("DELETE /api/monitors/{id}", h.HandleDelete)
}
