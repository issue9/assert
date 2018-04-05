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
)

// Request 请求的参数封装
type Request struct {
	path    string
	method  string
	body    io.Reader
	queries url.Values
	params  map[string]string

	srv *Server
}

// NewRequest 获取一条请求的结果
func (srv *Server) NewRequest(method, path string) *Request {
	return &Request{
		path:   path,
		method: method,
		srv:    srv,
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

// Body 指定提交的内容
func (req *Request) Body(body []byte) *Request {
	req.body = bytes.NewReader(body)
	return req
}

// JSONBody 指定一个 JSON 格式的 body
func (req *Request) JSONBody(obj interface{}) *Request {
	data, err := json.Marshal(obj)
	req.srv.assertion.NotError(err).NotNil(data)
	return req.Body(data)
}

// Param 替换参数
func (req *Request) Param(key, val string) *Request {
	if req.params == nil {
		req.params = make(map[string]string, 5)
	}

	req.params[key] = val

	return req
}

// Do 执行请求操作
func (req *Request) Do() *Response {
	path := req.path
	for key, val := range req.params {
		key = "{" + key + "}"
		path = strings.Replace(path, key, val, -1)
	}

	r, err := http.NewRequest(req.method, req.srv.server.URL+path, req.body)
	req.srv.assertion.NotError(err).NotNil(r)

	resp, err := req.srv.client.Do(r)
	req.srv.assertion.NotError(err).NotNil(resp)

	return &Response{
		assertion: req.srv.assertion,
		resp:      resp,
	}
}
