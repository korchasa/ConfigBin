package server_test

import (
	"configBin/pkg/encryptor/aes"
	"configBin/pkg/metrics/fake"
	"configBin/pkg/server"
	"configBin/pkg/server/responder"
	"configBin/pkg/server/templates"
	"configBin/pkg/server/utils"
	"configBin/pkg/storage/sqlite"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func NewTestingServer(sqlitePath string) (*server.Server, *sqlite.Storage, error) {
	enc := aes.NewEncryptor()

	err := os.Remove(sqlitePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("failed to remove sqlite file: %w", err)
	}

	store, err := sqlite.NewSqliteStorage(sqlitePath, enc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create sqlite storage: %w", err)
	}
	err = store.ApplySchema()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply db schema: %w", err)
	}

	tplProvider, err := templates.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build templates: %w", err)
	}

	resp := responder.New(tplProvider)
	srv := server.New(store, resp, tplProvider, fake.Fake{})

	return srv, store, nil
}

type formRequestSpec struct {
	method         string
	path           string
	formData       string
	cookieBinID    uuid.UUID
	cookiePassword string
}

func requestWithForm(spec formRequestSpec) *http.Request {
	form := strings.NewReader(spec.formData)
	req, _ := http.NewRequest(spec.method, spec.path, form)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func requestWithFormAndCookie(spec formRequestSpec) *http.Request {
	req := requestWithForm(spec)
	req.AddCookie(utils.PasswordCookie(spec.cookieBinID, spec.cookiePassword))
	return req
}

func createBinForTest(t *testing.T, store *sqlite.Storage) uuid.UUID {
	t.Helper()
	binID := uuid.New()
	err := store.CreateBin(binID, "test1234", "test_content")
	assert.NoError(t, err)
	return binID
}
