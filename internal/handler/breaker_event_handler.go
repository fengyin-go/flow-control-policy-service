package handler

import (
	"net/http"

	"flowcontrol/pkg/httpx"
)

func (s *Server) registerBreakerEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/breaker-events", s.listBreakerEvents)
}

func (s *Server) listBreakerEvents(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListBreakerEvents(r.URL.Query().Get("breaker_id")))
}
