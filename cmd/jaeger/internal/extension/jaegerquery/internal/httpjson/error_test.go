// Copyright (c) 2026 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package httpjson

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaegertracing/jaeger/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.VerifyGoLeaks(m)
}

func TestWriteError(t *testing.T) {
	testCases := []struct {
		name         string
		response     any
		statusCode   int
		expectedBody string
	}{
		{
			name:         "struct response",
			response:     map[string]string{"error": "bad request"},
			statusCode:   http.StatusBadRequest,
			expectedBody: "{\"error\":\"bad request\"}\n",
		},
		{
			name:         "bytes response without newline",
			response:     []byte(`{"msg":"not found"}`),
			statusCode:   http.StatusNotFound,
			expectedBody: "{\"msg\":\"not found\"}\n",
		},
		{
			name:         "bytes response with newline",
			response:     []byte("{\"msg\":\"server error\"}\n"),
			statusCode:   http.StatusInternalServerError,
			expectedBody: "{\"msg\":\"server error\"}\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.Header().Set("Content-Length", "100")

			WriteError(w, tc.statusCode, tc.response)

			assert.Empty(t, w.Header().Get("Content-Length"))
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
			assert.Equal(t, tc.statusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestWriteErrorMarshalFailure(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set("Content-Length", "50")

	// Channels cannot be marshaled to JSON.
	WriteError(w, http.StatusBadRequest, make(chan int))

	assert.Empty(t, w.Header().Get("Content-Length"))
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"errors":[{"code":500,"msg":"failed to marshal error response"}]}`, w.Body.String())
}
