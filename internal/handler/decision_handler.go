package handler

import (
	"net/http"

	"flowcontrol/pkg/httpx"
)

func (s *Server) registerDecisionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/decide", s.decide)
	mux.HandleFunc("POST /api/records/success", s.recordSuccess)
	mux.HandleFunc("POST /api/records/failure", s.recordFailure)
}

type decideRequest struct {
	Resource string `json:"resource"`
	Key      string `json:"key"`
}

func (s *Server) decide(w http.ResponseWriter, r *http.Request) {
	var req decideRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	httpx.OK(w, s.svc.Decide(req.Resource, req.Key))
}

type recordRequest struct {
	Resource string `json:"resource"`
}

func (s *Server) recordSuccess(w http.ResponseWriter, r *http.Request) {
	var req recordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.RecordSuccess(req.Resource); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"status": "recorded"})
}

func (s *Server) recordFailure(w http.ResponseWriter, r *http.Request) {
	var req recordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.RecordFailure(req.Resource); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"status": "recorded"})
}
