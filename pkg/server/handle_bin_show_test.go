package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBinShow(t *testing.T) {
	srv, store, err := NewTestingServer("./TestHandleBinShow.sqlite")
	require.NoError(t, err)
	t.Parallel()

	// Test case: Successful bin show
	t.Run("successful bin show", func(t *testing.T) {
		t.Parallel()

		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "GET",
			path:           "/" + binID.String(),
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "test_content")
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		t.Parallel()

		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "GET",
			path:           "/invalid",
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		t.Parallel()

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "GET",
			path:           "/00000000-0000-0000-0000-000000000000",
			cookieBinID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_get_bin")
	})

	// Test case: Missing password cookie
	t.Run("missing password cookie", func(t *testing.T) {
		t.Parallel()

		binID := createBinForTest(t, store)

		req := requestWithFormAndCookie(formRequestSpec{
			method:         "GET",
			path:           "/" + binID.String(),
			cookieBinID:    binID,
			cookiePassword: "",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Enter the bin password")
	})
}
