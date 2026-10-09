// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v5"
)

func TestServer(t *testing.T) {
	a := assert.New(t, false)

	testServer(NewTestServer(a, h))
	testServer(NewServer(a, h))
	testServer(NewTLSServer(a, h))
}

func testServer(s *Server) {
	s.Get("/get").Do().Status(http.StatusCreated)
	s.Post("/body", nil).
		Header("content-type", "application/json").
		StringBody(`{"id":5}`).
		Do().
		Status(http.StatusCreated).
		StringBody(`{"id":6}`)

	s.Get("/not-exists").Do().Status(http.StatusNotFound)

	s.Close()
	s.Close()
}
