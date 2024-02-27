// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestRequest_Do(t *testing.T) {
	a := assert.New(t, false)
	srv := NewServer(a, h, nil)

	srv.Get("/get").
		Do(nil).
		Success().
		Status(201)

	srv.NewRequest(http.MethodGet, "/not-exists").
		Do(nil).
		Fail()

	srv.NewRequest(http.MethodGet, "/get").
		Do(BuildHandler(a, 202, "", nil)).
		Status(202)

	r := Get(a, "/get")
	r.Do(BuildHandler(a, 202, "", nil)).Status(202)
	r.Do(BuildHandler(a, 203, "", nil)).Status(203)
	a.Panic(func() {
		r.Do(nil)
	})
}

func TestResponse(t *testing.T) {
	srv := NewServer(assert.New(t, false), h, nil)

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/json").
		Query("page", "5").
		JSONBody(&bodyTest{ID: 5}).
		Do(nil).
		Status(http.StatusCreated).
		NotStatus(http.StatusNotFound).
		Header("content-type", "application/json;charset=utf-8").
		NotHeader("content-type", "invalid value").
		JSONBody(&bodyTest{ID: 6}).
		Body([]byte(`{"id":6}`)).
		StringBody(`{"id":6}`).
		BodyNotEmpty()

	srv.NewRequest(http.MethodGet, "/get").
		Query("page", "5").
		Do(nil).
		Status(http.StatusCreated).
		NotHeader("content-type", "invalid value").
		BodyEmpty()

	// xml

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/xml").
		XMLBody(&bodyTest{ID: 5}).
		Do(nil).
		Success(http.StatusCreated).
		Header("content-type", "application/xml;charset=utf-8").
		XMLBody(&bodyTest{ID: 6})
}
