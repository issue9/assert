// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"net/http/httptest"

	"github.com/issue9/assert"
)

// Server 测试服务
type Server struct {
	a      *assert.Assertion
	server *httptest.Server
	client *http.Client
}

// NewServer 声明新的测试服务
//
// 如果 client 为 nil，则会采用 http.DefaultClient 作为默认值
func NewServer(a *assert.Assertion, h http.Handler, client *http.Client) *Server {
	return newServer(a, httptest.NewServer(h), client)
}

// NewTLSServer 声明新的测试服务
//
// 如果 client 为 nil，则会采用 http.DefaultClient 作为默认值
func NewTLSServer(a *assert.Assertion, h http.Handler, client *http.Client) *Server {
	return newServer(a, httptest.NewTLSServer(h), client)
}

func newServer(a *assert.Assertion, srv *httptest.Server, client *http.Client) *Server {
	if client == nil {
		client = &http.Client{}
	}

	return &Server{
		a:      a,
		server: srv,
		client: client,
	}
}

// Close 停止服务
func (srv *Server) Close() { srv.server.Close() }

func (srv *Server) Assertion() *assert.Assertion { return srv.a }
