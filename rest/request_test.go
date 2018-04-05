// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

var h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/get" {
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusNotFound)
})

func TestDefaultRequest(t *testing.T) {
	a := assert.New(t)

	srv := NewServer(a, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	req := srv.NewRequest(http.MethodGet, "/get")
	a.NotNil(req)
	resp := req.Do()
	resp.Status(http.StatusCreated)

	req = srv.NewRequest(http.MethodGet, "/not-exists")
	a.NotNil(req)
	resp = req.Do()
	resp.Status(http.StatusNotFound)
}
