package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBinUpdate(t *testing.T) {
	srv, store, err := NewTestingServer("./TestHandleBinUpdate.sqlite")
	require.NoError(t, err)
	t.Parallel()

	// Test case: Successful bin update
	t.Run("successful bin update", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/" + binID.String() + "/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test1234",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusFound, resp.Code)
		log.Info(resp.Body.String())
		assert.Equal(t, "/"+binID.String(), resp.Header().Get("Location"))
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/invalid/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test1234",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/00000000-0000-0000-0000-000000000000/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test1234",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_update_bin")
	})

	// Test case: Missing content
	t.Run("missing content", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/" + binID.String() + "/update",
			formData:       "",
			cookieBinID:    binID,
			cookiePassword: "test1234",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "no_content")
	})

	// Test case: Missing password
	t.Run("missing cookie", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:   "POST",
			path:     "/" + binID.String() + "/update",
			formData: "content=updated_content",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_update_bin")
	})

	// Test case: Wrong password
	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		binID := createBinForTest(t, store)

		req := requestWithForm(formRequestSpec{
			method:         "POST",
			path:           "/" + binID.String() + "/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "wrong_password",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_update_bin")
	})
}
