// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v5"
)

func TestRequest_Do(t *testing.T) {
	a := assert.New(t, false)

	Get(a, "/get", h).
		Do().
		Success().
		Status(201)

	NewRequest(a, http.MethodGet, "/not-exists", h).
		Do().
		Fail()

	NewRequest(a, http.MethodGet, "/get", BuildHandler(a, 202, "", nil)).
		Do().
		Status(202)

	Get(a, "/get", BuildHandler(a, 202, "", nil)).Do().Status(202)
	Get(a, "/get", BuildHandler(a, 203, "", nil)).Do().Status(203)
	a.Panic(func() {
		Get(a, "/get", nil).Do()
	})
}

func TestResponse(t *testing.T) {
	srv := NewServer(assert.New(t, false), h)

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/json").
		Query("page", "5").
		StringBody(`{"id":5}`).
		Do().
		Status(http.StatusCreated).
		NotStatus(http.StatusNotFound).
		Header("content-type", "application/json;charset=utf-8").
		NotHeader("content-type", "invalid value").
		Body([]byte(`{"id":6}`)).
		StringBody(`{"id":6}`).
		BodyNotEmpty()

	srv.NewRequest(http.MethodGet, "/get").
		Query("page", "5").
		Do().
		Status(http.StatusCreated).
		NotHeader("content-type", "invalid value").
		BodyEmpty()
}
