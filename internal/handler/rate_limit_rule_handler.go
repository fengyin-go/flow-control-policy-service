package handler

import (
	"net/http"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createRule)
	mux.HandleFunc("GET /api/rules", s.listRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateRule)
	mux.HandleFunc("PATCH /api/rules/{id}/toggle", s.toggleRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRule)
}

type createRuleRequest struct {
	Name      string `json:"name"`
	Resource  string `json:"resource"`
	Algorithm string `json:"algorithm"`
	Limit     int    `json:"limit"`
	WindowSec int    `json:"window_sec"`
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateRateLimitRule(model.RateLimitRule{
		Name:      req.Name,
		Resource:  req.Resource,
		Algorithm: req.Algorithm,
		Limit:     req.Limit,
		WindowSec: req.WindowSec,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleFilter{
		Resource:  r.URL.Query().Get("resource"),
		Algorithm: r.URL.Query().Get("algorithm"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRateLimitRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.svc.GetRateLimitRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateRateLimitRule(r.PathValue("id"), model.RateLimitRule{
		Name:      req.Name,
		Resource:  req.Resource,
		Algorithm: req.Algorithm,
		Limit:     req.Limit,
		WindowSec: req.WindowSec,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

type toggleRuleRequest struct {
	Enabled bool `json:"enabled"`
}

func (s *Server) toggleRule(w http.ResponseWriter, r *http.Request) {
	var req toggleRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.ToggleRateLimitRule(r.PathValue("id"), req.Enabled)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRateLimitRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
