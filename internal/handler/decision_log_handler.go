package handler

import (
	"net/http"
	"strconv"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerDecisionLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/decision-logs", s.listDecisionLogs)
	mux.HandleFunc("GET /api/decision-stats", s.decisionStats)
}

func (s *Server) listDecisionLogs(w http.ResponseWriter, r *http.Request) {
	filter := model.DecisionLogFilter{Resource: r.URL.Query().Get("resource")}
	if v := r.URL.Query().Get("allowed"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.Allowed = &b
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	httpx.OK(w, s.svc.ListDecisionLogs(filter, limit))
}

func (s *Server) decisionStats(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.DecisionStats(r.URL.Query().Get("resource")))
}
