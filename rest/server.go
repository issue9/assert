// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package rest

import (
	"net/http"
	"net/http/httptest"

	"github.com/issue9/assert/v5"
)

// Server 测试服务
type Server struct {
	a      *assert.Assertion
	server *httptest.Server
}

// NewServer 声明新的测试服务
//
// 这是使用 [httptest.NewTestServer] 创建的内存服务，无须端口。
func NewServer(a *assert.Assertion, h http.Handler) *Server {
	return &Server{
		a:      a,
		server: httptest.NewTestServer(a.TB(), h),
	}
}

// URL 测试服务的基地址
//
// 该值可能不是一个固定的值，具体可参考 [httptest.Server.URL] 的相关说明。
func (srv *Server) URL() string {
	if srv.server.URL != "" {
		return srv.server.URL
	}
	return "https://example.com"
}

func (srv *Server) Assertion() *assert.Assertion { return srv.a }

// Close 关闭服务
//
// 该方法会自动注册在 [testing.T.Cleanup]，一般情况下无需手动调用。
func (srv *Server) Close() { srv.server.Close() }

// Server 返回 [httptest.Server] 对象
func (srv *Server) Server() *httptest.Server { return srv.server }
