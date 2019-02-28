// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

func TestRequest_Do(t *testing.T) {
	srv := NewServer(t, h, nil)
	assert.NotNil(t, srv)
	defer srv.Close()

	srv.NewRequest(http.MethodGet, "/get").
		Do().
		Success()

	srv.NewRequest(http.MethodGet, "/not-exists").
		Do().
		Fail()
}

func TestResponse(t *testing.T) {
	srv := NewServer(t, h, nil)
	assert.NotNil(t, srv)
	defer srv.Close()

	w1 := new(bytes.Buffer)
	w2 := new(bytes.Buffer)

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/json").
		Query("page", "5").
		JSONBody(&bodyTest{ID: 5}).
		Do().
		Status(http.StatusCreated).
		NotStatus(http.StatusNotFound).
		Header("content-type", "application/json;charset=utf-8").
		NotHeader("content-type", "invalid value").
		ReadBody(w1).
		ReadBody(w2).
		JSONBody(&bodyTest{ID: 6}).
		BodyNotNil().
		BodyNotEmpty()
	assert.Equal(t, w1.String(), w2.String())

	srv.NewRequest(http.MethodGet, "/get").
		Query("page", "5").
		Do().
		Status(http.StatusCreated).
		NotHeader("content-type", "invalid value").
		BodyNotNil().
		BodyEmpty()
}
