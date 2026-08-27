package handler

import (
	"net/http"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alert-rules", s.createAlertRule)
	mux.HandleFunc("GET /api/alert-rules", s.listAlertRules)
	mux.HandleFunc("GET /api/alert-rules/evaluate", s.evaluateAlerts)
	mux.HandleFunc("GET /api/alert-rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/alert-rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/alert-rules/{id}", s.deleteAlertRule)
}

type createAlertRuleRequest struct {
	Name      string  `json:"name"`
	Resource  string  `json:"resource"`
	Metric    string  `json:"metric"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAlertRule(model.AlertRule{
		Name:      req.Name,
		Resource:  req.Resource,
		Metric:    req.Metric,
		Threshold: req.Threshold,
		Operator:  req.Operator,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListAlertRules())
}

func (s *Server) evaluateAlerts(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.EvaluateAlerts())
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	a, err := s.svc.GetAlertRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAlertRule(r.PathValue("id"), model.AlertRule{
		Name:      req.Name,
		Resource:  req.Resource,
		Metric:    req.Metric,
		Threshold: req.Threshold,
		Operator:  req.Operator,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAlertRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
