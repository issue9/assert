// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestRequest_Do(t *testing.T) {
	srv := NewServer(assert.New(t, false), h, nil)
	defer srv.Close()

	srv.NewRequest(http.MethodGet, "/get").
		Do().
		Success()

	srv.NewRequest(http.MethodGet, "/not-exists").
		Do().
		Fail()
}

func TestResponse(t *testing.T) {
	srv := NewServer(assert.New(t, false), h, nil)
	defer srv.Close()

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/json").
		Query("page", "5").
		JSONBody(&bodyTest{ID: 5}).
		Do().
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
		Do().
		Status(http.StatusCreated).
		NotHeader("content-type", "invalid value").
		BodyEmpty()

	// xml

	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/xml").
		XMLBody(&bodyTest{ID: 5}).
		Do().
		Success(http.StatusCreated).
		Header("content-type", "application/xml;charset=utf-8").
		XMLBody(&bodyTest{ID: 6})
}
