// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/issue9/assert"
)

// Request 请求的参数封装
type Request struct {
	path      string
	method    string
	body      io.Reader
	queries   url.Values
	params    map[string]string
	headers   map[string]string
	assertion *assert.Assertion

	client *http.Client
	prefix string // 地址前缀
}

// NewRequest 获取一条请求的结果
func (srv *Server) NewRequest(method, path string) *Request {
	req := NewRequest(srv.assertion, srv.client, method, path)
	req.prefix = srv.server.URL

	return req
}

// NewRequest 创建新的请求实例
//
// client 如果为空，则会采用 http.DefaultClient 作为其值。
func NewRequest(a *assert.Assertion, client *http.Client, method, path string) *Request {
	if client == nil {
		client = http.DefaultClient
	}

	return &Request{
		assertion: a,
		client:    client,
		method:    method,
		path:      path,
	}
}

// Query 替换一个请求参数
func (req *Request) Query(key, val string) *Request {
	if req.queries == nil {
		req.queries = url.Values{}
	}

	req.queries.Add(key, val)

	return req
}

// Header 指定请求时的报头
func (req *Request) Header(key, val string) *Request {
	if req.headers == nil {
		req.headers = make(map[string]string, 5)
	}

	req.headers[key] = val

	return req
}

// Param 替换参数
func (req *Request) Param(key, val string) *Request {
	if req.params == nil {
		req.params = make(map[string]string, 5)
	}

	req.params[key] = val

	return req
}

// Body 指定提交的内容
func (req *Request) Body(body []byte) *Request {
	req.body = bytes.NewReader(body)
	return req
}

// JSONBody 指定一个 JSON 格式的 body
func (req *Request) JSONBody(obj interface{}) *Request {
	data, err := json.Marshal(obj)
	req.assertion.NotError(err).NotNil(data)
	return req.Body(data)
}

func (req *Request) buildPath() string {
	path := req.path

	for key, val := range req.params {
		key = "{" + key + "}"
		path = strings.Replace(path, key, val, -1)
	}

	if len(req.queries) > 0 {
		path += ("?" + req.queries.Encode())
	}

	return path
}
