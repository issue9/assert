// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

func TestNew(t *testing.T) {
	a := assert.New(t)

	srv := NewTLSServer(t, nil, nil)
	a.NotNil(srv)
	a.Equal(srv.client, http.DefaultClient)
	a.True(len(srv.server.URL) > 0)
	srv.Close()

	client := &http.Client{}
	srv = NewServer(t, nil, client)
	a.NotNil(srv)
	a.Equal(srv.client, client)
	srv.Close()
}
