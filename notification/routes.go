package notification

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /api/monitors/{id}/notifications", h.HandleList)
	mux.HandleFunc("POST /api/monitors/{id}/notifications", h.HandleCreate)
	mux.HandleFunc("GET /api/notifications/{nid}", h.HandleGet)
	mux.HandleFunc("PUT /api/notifications/{nid}", h.HandleUpdate)
	mux.HandleFunc("DELETE /api/notifications/{nid}", h.HandleDelete)
}
