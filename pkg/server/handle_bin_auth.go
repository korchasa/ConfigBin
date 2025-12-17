package server

import (
	"configBin/pkg/server/utils"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var ErrBinNotFoundByID = errors.New("no bin found by id")

// handleShowEventSchema is the handler to show event schema.
func (s *Server) handleBinAuth() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		bid, err := utils.ExtractBinIDFromPathVar(req)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		err = req.ParseForm()
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_parse_form", err)
			return
		}

		pass := req.Form.Get("password")
		if pass == "" {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "empty_password", ErrEmptyPassword)
			return
		}

		bin, err := s.store.GetBin(*bid, pass)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_get_bin", err)
			return
		}
		if bin == nil {
			s.resp.HTMLError(req, resp, http.StatusNotFound, "no_bin_by_id", ErrBinNotFoundByID)
			return
		}

		http.SetCookie(resp, utils.PasswordCookie(*bid, pass))

		log.Infof("bin authed: %s", bid.String())

		http.Redirect(resp, req, "/"+bid.String(), http.StatusFound)
	}
}
