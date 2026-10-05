// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"testing"

	"github.com/issue9/assert/v5"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	srv := NewServer(a, nil)
	a.NotNil(srv).NotNil(srv.Assertion())

	srv.Close()
	srv.Close()
}
