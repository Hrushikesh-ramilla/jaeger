// Copyright (c) 2026 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

// Package httpjson provides JSON response helpers for Jaeger Query HTTP APIs.
package httpjson

import (
	"encoding/json"
	"net/http"
)

// WriteError writes an error response with a JSON body and appropriate headers.
// It sets Content-Type to application/json, removes Content-Length, sets
// X-Content-Type-Options: nosniff, and writes the given status code.
// If marshaling fails, it writes an internal server error JSON response.
func WriteError(w http.ResponseWriter, statusCode int, response any) {
	var body []byte
	switch r := response.(type) {
	case []byte:
		body = r
	default:
		var err error
		body, err = json.Marshal(response)
		if err != nil {
			w.Header().Del("Content-Length")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("{\"errors\":[{\"code\":500,\"msg\":\"failed to marshal error response\"}]}\n"))
			return
		}
	}

	w.Header().Del("Content-Length")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		body = append(body, '\n')
	}
	_, _ = w.Write(body)
}
