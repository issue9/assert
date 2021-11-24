// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestRequest_buildPath(t *testing.T) {
	srv := NewServer(assert.New(t, false), h, nil)
	a := srv.Assertion()
	a.NotNil(srv)

	req := srv.NewRequest(http.MethodGet, "/get")
	a.NotNil(req)
	a.Equal(req.buildPath(), srv.server.URL+"/get")

	req.Param("id", "1").Query("page", "5")
	a.Equal(req.buildPath(), srv.server.URL+"/get?page=5")

	req = srv.NewRequest(http.MethodGet, "/users/{id}/orders/{oid}")
	a.NotNil(req)
	a.Equal(req.buildPath(), srv.server.URL+"/users/{id}/orders/{oid}")
	req.Param("id", "1").Param("oid", "2").Query("page", "5")
	a.Equal(req.buildPath(), srv.server.URL+"/users/1/orders/2?page=5")
}
