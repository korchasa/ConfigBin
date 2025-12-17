package server

import "net/http"

// handleReadiness is the handler to check if the service is ready to serve requests.
func (s *Server) handleReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if s.store.IsReady() {
			s.resp.JSONSuccess(w, http.StatusOK, "readiness:OK")
		} else {
			s.resp.JSONSuccess(w, http.StatusServiceUnavailable, "readiness:FAIL")
		}
	}
}
