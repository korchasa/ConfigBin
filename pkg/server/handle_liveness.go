package server

import "net/http"

// handleLiveness is the handler to check if the service is alive.
func (s *Server) handleLiveness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		s.resp.JSONSuccess(w, http.StatusOK, "liveness:OK")
	}
}
