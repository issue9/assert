// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestBuildHandler(t *testing.T) {
	a := assert.New(t, false)

	h := BuildHandler(a, 201, "body", map[string]string{"k1": "v1"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)
	a.Equal(w.Code, 201)
}
