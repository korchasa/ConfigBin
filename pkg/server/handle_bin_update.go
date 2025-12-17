package server

import (
	"configBin/pkg/server/utils"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// handleBinUpdate is the handler to update a bin.
func (s *Server) handleBinUpdate() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		binID, err := utils.ExtractBinIDFromPathVar(req)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "invalid_bin_id", err)
			return
		}

		pass := utils.ReadPasswordCookie(req, *binID)

		err = req.ParseForm()
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_parse_form", err)
			return
		}

		unencryptedData := req.Form.Get("content")
		if unencryptedData == "" {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "no_content", ErrEmptyContent)
			return
		}

		err = s.store.UpdateBin(*binID, pass, unencryptedData)
		if err != nil {
			s.resp.HTMLError(req, resp, http.StatusBadRequest, "cant_update_bin", err)
			return
		}

		log.Infof("bin updated: %s", binID.String())

		http.Redirect(resp, req, "/"+binID.String(), http.StatusFound)
	}
}
