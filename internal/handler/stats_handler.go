package handler

import (
	"net/http"

	"flowcontrol/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/flow", s.statsFlow)
	mux.HandleFunc("GET /api/stats/resource", s.resourceStats)
}

func (s *Server) statsFlow(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsFlow())
}

func (s *Server) resourceStats(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ResourceStats(r.URL.Query().Get("resource")))
}
