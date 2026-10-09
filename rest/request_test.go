// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/issue9/assert/v5"
)

func TestRequest_buildPath(t *testing.T) {
	srv := NewTestServer(assert.New(t, false), h)
	a := srv.Assertion()
	a.NotNil(srv)

	req := srv.NewRequest(http.MethodGet, "/get")
	a.NotNil(req)
	u, err := url.Parse(req.buildPath())
	a.NotError(err).Equal(u.Path, "/get")

	req.Param("id", "1").Query("page", "5")
	u, err = url.Parse(req.buildPath())
	a.NotError(err).Equal(u.Path+"?"+u.RawQuery, "/get?page=5")

	req = srv.NewRequest(http.MethodGet, "/users/{id}/orders/{oid}")
	a.NotNil(req)
	u, err = url.Parse(req.buildPath())
	a.NotError(err).Equal(u.Path, "/users/{id}/orders/{oid}")
	req.Param("id", "1").Param("oid", "2").Query("page", "5")
	u, err = url.Parse(req.buildPath())
	a.NotError(err).Equal(u.Path+"?"+u.RawQuery, "/users/1/orders/2?page=5")
}
