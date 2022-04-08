// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"net/http/httptest"

	"github.com/issue9/assert/v2"
)

// Server 测试服务
type Server struct {
	a      *assert.Assertion
	server *httptest.Server
	client *http.Client
	closed bool
}

// NewServer 声明新的测试服务
//
// 如果 client 为 nil，则会采用 &http.Client{} 作为默认值
func NewServer(a *assert.Assertion, h http.Handler, client *http.Client) *Server {
	return newServer(a, httptest.NewServer(h), client)
}

// NewTLSServer 声明新的测试服务
//
// 如果 client 为 nil，则会采用 &http.Client{} 作为默认值
func NewTLSServer(a *assert.Assertion, h http.Handler, client *http.Client) *Server {
	return newServer(a, httptest.NewTLSServer(h), client)
}

func newServer(a *assert.Assertion, srv *httptest.Server, client *http.Client) *Server {
	if client == nil {
		client = &http.Client{}
	}

	s := &Server{
		a:      a,
		server: srv,
		client: client,
	}

	a.TB().Cleanup(func() {
		s.Close()
	})

	return s
}

func (srv *Server) URL() string { return srv.server.URL }

func (srv *Server) Assertion() *assert.Assertion { return srv.a }

// Close 关闭服务
//
// 如果未手动调用，则在 testing.TB.Cleanup 中自动调用。
func (srv *Server) Close() {
	if srv.closed {
		return
	}

	srv.server.Close()
	srv.closed = true
}
