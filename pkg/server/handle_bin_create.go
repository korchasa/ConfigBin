package server

import (
	"configBin/pkg/server/utils"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// handleBinCreate is the handler to create a bin.
func (s *Server) handleBinCreate() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_parse_form", err)
			return
		}

		id := req.Form.Get("uuid")
		if id == "" {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "empty_uuid", ErrEmptyUUID)
			return
		}
		bid, err := uuid.Parse(id)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "invalid_bin_id", ErrInvalidBinID)
			return
		}

		pass := req.Form.Get("password")
		if pass == "" {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "empty_password", ErrEmptyPassword)
			return
		}
		if len(pass) < 8 {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "password_too_short", ErrPasswordTooShort)
			return
		}

		unencryptedData := req.Form.Get("content")
		if unencryptedData == "" {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "empty_content", ErrEmptyContent)
		}

		err = s.store.CreateBin(bid, pass, unencryptedData)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_create_bin", err)
			return
		}

		http.SetCookie(resp, utils.PasswordCookie(bid, pass))

		log.Infof("bin created: %s", bid.String())

		http.Redirect(resp, req, "/"+bid.String(), http.StatusFound)
	}
}
