package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleBinCreate(t *testing.T) {
	srv, _, err := NewTestingServer("./TestHandleBinCreate.sqlite")
	require.NoError(t, err)
	t.Parallel()

	// Test case: Successful bin creation
	t.Run("successful bin creation", func(t *testing.T) {
		t.Parallel()

		bid := uuid.New().String()
		req := requestWithForm(formRequestSpec{
			method:   "POST",
			path:     "/create",
			formData: "uuid=" + bid + "&password=test&content=test_content",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusFound, resp.Code)
		assert.Equal(t, "/"+bid, resp.Header().Get("Location"))
	})

	// Test case: Missing UUID
	t.Run("missing uuid", func(t *testing.T) {
		t.Parallel()

		req := requestWithForm(formRequestSpec{
			method:   "POST",
			path:     "/create",
			formData: "password=test&content=test_content",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "empty_uuid")
	})

	// Test case: Invalid UUID
	t.Run("invalid uuid", func(t *testing.T) {
		t.Parallel()

		req := requestWithForm(formRequestSpec{
			method:   "POST",
			path:     "/create",
			formData: "uuid=invalid&password=test&content=test_content",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Missing password
	t.Run("missing password", func(t *testing.T) {
		t.Parallel()

		bid := uuid.New().String()
		req := requestWithForm(formRequestSpec{
			method:   "POST",
			path:     "/create",
			formData: "uuid=" + bid + "&content=test_content",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "empty_password")
	})

	// Test case: Missing content
	t.Run("missing content", func(t *testing.T) {
		t.Parallel()

		bid := uuid.New().String()
		req := requestWithForm(formRequestSpec{
			method:   "POST",
			path:     "/create",
			formData: "uuid=" + bid + "&password=test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "empty_content")
	})
}
