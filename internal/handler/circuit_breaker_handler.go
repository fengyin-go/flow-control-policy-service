package handler

import (
	"net/http"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerBreakerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/breakers", s.createBreaker)
	mux.HandleFunc("GET /api/breakers", s.listBreakers)
	mux.HandleFunc("GET /api/breakers/{id}", s.getBreaker)
	mux.HandleFunc("PUT /api/breakers/{id}", s.updateBreaker)
	mux.HandleFunc("DELETE /api/breakers/{id}", s.deleteBreaker)
}

type createBreakerRequest struct {
	Name             string `json:"name"`
	Resource         string `json:"resource"`
	FailureThreshold int    `json:"failure_threshold"`
	SuccessThreshold int    `json:"success_threshold"`
	TimeoutMs        int    `json:"timeout_ms"`
	OpenDurationSec  int    `json:"open_duration_sec"`
}

func (s *Server) createBreaker(w http.ResponseWriter, r *http.Request) {
	var req createBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateCircuitBreaker(model.CircuitBreaker{
		Name:             req.Name,
		Resource:         req.Resource,
		FailureThreshold: req.FailureThreshold,
		SuccessThreshold: req.SuccessThreshold,
		TimeoutMs:        req.TimeoutMs,
		OpenDurationSec:  req.OpenDurationSec,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBreakers(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListCircuitBreakers())
}

func (s *Server) getBreaker(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetCircuitBreaker(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) updateBreaker(w http.ResponseWriter, r *http.Request) {
	var req createBreakerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateCircuitBreaker(r.PathValue("id"), model.CircuitBreaker{
		Name:             req.Name,
		FailureThreshold: req.FailureThreshold,
		SuccessThreshold: req.SuccessThreshold,
		TimeoutMs:        req.TimeoutMs,
		OpenDurationSec:  req.OpenDurationSec,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBreaker(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteCircuitBreaker(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
