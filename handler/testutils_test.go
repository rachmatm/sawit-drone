package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

// newEchoContext creates an Echo context with an optional JSON body.
// Use it in every handler test instead of setting up Echo manually each time.
func newEchoContext(t *testing.T, method, path string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	var reqBody *bytes.Reader
	switch v := body.(type) {
	case nil:
		reqBody = bytes.NewReader(nil)
	case string:
		// use raw string as-is, no marshalling
		reqBody = bytes.NewReader([]byte(v))
	default:
		b, err := json.Marshal(v)
		require.NoError(t, err)
		reqBody = bytes.NewReader(b)
	}

	e := echo.New()
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec), rec
}

// helper: decode response body
func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	require.NoError(t, json.NewDecoder(rec.Body).Decode(v))
}
