// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	srv := NewTLSServer(a, nil, nil)
	a.NotNil(srv)
	a.Equal(srv.client, &http.Client{})
	a.True(len(srv.server.URL) > 0)

	client := &http.Client{}
	srv = NewServer(a, nil, client)
	a.NotNil(srv)
	a.Equal(srv.client, client)

	srv.Close()
	a.True(srv.closed)
	srv.Close()
}
