package handler

import (
	"net/http"

	"flowcontrol/internal/model"
	"flowcontrol/pkg/httpx"
)

func (s *Server) registerTokenBucketRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/token-buckets", s.createTokenBucket)
	mux.HandleFunc("GET /api/token-buckets", s.listTokenBuckets)
	mux.HandleFunc("GET /api/token-buckets/{id}", s.getTokenBucket)
	mux.HandleFunc("POST /api/token-buckets/{id}/refill", s.refillTokenBucket)
	mux.HandleFunc("POST /api/token-buckets/{id}/take", s.takeToken)
	mux.HandleFunc("DELETE /api/token-buckets/{id}", s.deleteTokenBucket)
}

type createTokenBucketRequest struct {
	RuleID     string  `json:"rule_id"`
	Capacity   int     `json:"capacity"`
	RefillRate int     `json:"refill_rate"`
	Tokens     float64 `json:"tokens"`
}

func (s *Server) createTokenBucket(w http.ResponseWriter, r *http.Request) {
	var req createTokenBucketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTokenBucket(model.TokenBucket{
		RuleID:     req.RuleID,
		Capacity:   req.Capacity,
		RefillRate: req.RefillRate,
		Tokens:     req.Tokens,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTokenBuckets(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.ListTokenBuckets())
}

func (s *Server) getTokenBucket(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTokenBucket(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) refillTokenBucket(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.RefillTokenBucket(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) takeToken(w http.ResponseWriter, r *http.Request) {
	ok, err := s.svc.TakeToken(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]bool{"taken": ok})
}

func (s *Server) deleteTokenBucket(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTokenBucket(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
