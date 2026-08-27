package handler

import (
	"net/http"

	"flowcontrol/pkg/httpx"
)

func (s *Server) registerTrafficStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/traffic-stats", s.listTrafficStats)
	mux.HandleFunc("GET /api/traffic-stats/{resource}/summary", s.trafficSummary)
}

func (s *Server) listTrafficStats(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListTrafficStats())
}

func (s *Server) trafficSummary(w http.ResponseWriter, r *http.Request) {
	sum, err := s.svc.TrafficSummary(r.PathValue("resource"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sum)
}
