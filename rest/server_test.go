// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	srv := NewTLSServer(a, nil, nil)
	a.NotNil(srv)
	a.Equal(srv.client, http.DefaultClient)
	a.True(len(srv.server.URL) > 0)
	srv.Close()

	client := &http.Client{}
	srv = NewServer(a, nil, client)
	a.NotNil(srv)
	a.Equal(srv.client, client)
	srv.Close()
}
