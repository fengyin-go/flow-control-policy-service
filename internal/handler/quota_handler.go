package handler

import (
	"net/http"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerQuotaRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/quotas", s.createQuota)
	mux.HandleFunc("GET /api/quotas", s.listQuotas)
	mux.HandleFunc("POST /api/quotas/consume", s.consumeQuota)
	mux.HandleFunc("DELETE /api/quotas/{id}", s.deleteQuota)
}

type createQuotaRequest struct {
	RuleID         string `json:"rule_id"`
	Dimension      string `json:"dimension"`
	DimensionValue string `json:"dimension_value"`
	Allowed        int    `json:"allowed"`
}

func (s *Server) createQuota(w http.ResponseWriter, r *http.Request) {
	var req createQuotaRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.CreateQuota(model.Quota{
		RuleID:         req.RuleID,
		Dimension:      req.Dimension,
		DimensionValue: req.DimensionValue,
		Allowed:        req.Allowed,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, q)
}

func (s *Server) listQuotas(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListQuotas(r.URL.Query().Get("rule_id")))
}

type consumeQuotaRequest struct {
	RuleID         string `json:"rule_id"`
	Dimension      string `json:"dimension"`
	DimensionValue string `json:"dimension_value"`
}

func (s *Server) consumeQuota(w http.ResponseWriter, r *http.Request) {
	var req consumeQuotaRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ok, err := s.svc.ConsumeQuota(req.RuleID, req.Dimension, req.DimensionValue)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]bool{"allowed": ok})
}

func (s *Server) deleteQuota(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteQuota(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
