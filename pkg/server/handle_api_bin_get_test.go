package server_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleAPIGetBin(t *testing.T) {
	srv, store, err := NewTestingServer("./TestHandleAPIGetBin.sqlite")
	require.NoError(t, err)
	t.Parallel()

	binID := uuid.New()
	binPass := "test1234:@\""
	binCredsHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "any-name", binPass)))
	err = store.CreateBin(binID, binPass, "test_content")
	require.NoError(t, err)

	// Test case: Successful bin retrieval
	t.Run("successful bin retrieval", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/"+binID.String(), nil)
		req.Header.Add("Authorization", "Basic "+binCredsHeader)
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "test_content")
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/invalid", nil)
		req.Header.Add("Authorization", "Basic "+binCredsHeader)
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Add("Authorization", "Basic "+binCredsHeader)
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_get_bin")
	})

	// Test case: Missing password
	t.Run("missing password", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/"+binID.String(), nil)
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "empty_password")
	})

	// Test case: Invalid password
	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/"+binID.String(), nil)
		req.Header.Add("Authorization", "Basic bm90OnZhbGlk")
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_get_bin")
	})
}
