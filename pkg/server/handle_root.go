package server

import (
	"configBin/pkg/server/utils"
	"net/http"

	"github.com/google/uuid"
)

const defaultPasswordLength = 8

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleRoot() http.HandlerFunc {
	tpl := s.tplProvider.MustGet("root.gohtml")
	return func(resp http.ResponseWriter, req *http.Request) {
		binID := uuid.New()
		password, err := utils.GeneratePassword(defaultPasswordLength)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusInternalServerError, "cant_generate_password", err)
		}
		err = tpl.Execute(resp, struct {
			ID       string
			Password string
		}{
			ID:       binID.String(),
			Password: password,
		})
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusInternalServerError, "cant_generate_password", err)
			return
		}
		resp.WriteHeader(http.StatusOK)
	}
}
