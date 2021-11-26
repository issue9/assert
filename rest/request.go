// SPDX-License-Identifier: MIT

package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/issue9/assert/v2"
)

// Request 请求的参数封装
type Request struct {
	path    string
	method  string
	body    io.Reader
	queries url.Values
	params  map[string]string
	headers map[string]string
	a       *assert.Assertion
	client  *http.Client
}

// NewRequest 获取一条请求的结果
//
// method 表示请求方法
// path 表示请求的路径，域名部分无须填定。可以通过 {} 指定参数，比如：
//  r := NewRequest(http.MethodGet, "/users/{id}")
// 之后就可以使用 Params 指定 id 的具体值，达到复用的目的：
//  resp1 := r.Param("id", "1").Do()
//  resp2 := r.Param("id", "2").Do()
func (srv *Server) NewRequest(method, path string) *Request {
	return &Request{
		a:      srv.a,
		client: srv.client,
		method: method,
		path:   srv.server.URL + path,
	}
}

// Get 相当于 NewRequest(http.MethodGet, path)
func (srv *Server) Get(path string) *Request {
	return srv.NewRequest(http.MethodGet, path)
}

// Put 相当于 NewRequest(http.MethodPut, path).Body()
func (srv *Server) Put(path string, body []byte) *Request {
	return srv.NewRequest(http.MethodPut, path).Body(body)
}

// Post 相当于 NewRequest(http.MethodPost, path).Body()
func (srv *Server) Post(path string, body []byte) *Request {
	return srv.NewRequest(http.MethodPost, path).Body(body)
}

// Patch 相当于 NewRequest(http.MethodPatch, path).Body()
func (srv *Server) Patch(path string, body []byte) *Request {
	return srv.NewRequest(http.MethodPatch, path).Body(body)
}

// Delete 相当于 NewRequest(http.MethodDelete, path).Body()
func (srv *Server) Delete(path string) *Request {
	return srv.NewRequest(http.MethodDelete, path)
}

// NewRequest 以调用链的方式构建一个访问请求对象
//
// NOTE: 通过此函数构建的对象只能访问 http.Handler 的内容，
// 如果需要访问网络内容，可以采用 Server.NewRequest。
func NewRequest(a *assert.Assertion, method, path string) *Request {
	return &Request{
		a:      a,
		method: method,
		path:   path,
	}
}

func Get(a *assert.Assertion, path string) *Request {
	return NewRequest(a, http.MethodGet, path)
}

func Delete(a *assert.Assertion, path string) *Request {
	return NewRequest(a, http.MethodDelete, path)
}

func Post(a *assert.Assertion, path string, body []byte) *Request {
	return NewRequest(a, http.MethodPost, path).Body(body)
}

func Put(a *assert.Assertion, path string, body []byte) *Request {
	return NewRequest(a, http.MethodPut, path).Body(body)
}

func Patch(a *assert.Assertion, path string, body []byte) *Request {
	return NewRequest(a, http.MethodPatch, path).Body(body)
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

func (req *Request) StringBody(body string) *Request {
	req.body = bytes.NewBufferString(body)
	return req
}

// BodyFunc 指定一个未编码的对象
//
// marshal 对 obj 的编码函数，比如 json.Marshal 等。
func (req *Request) BodyFunc(obj interface{}, marshal func(interface{}) ([]byte, error)) *Request {
	req.a.TB().Helper()

	data, err := marshal(obj)
	req.a.NotError(err).NotNil(data)
	return req.Body(data)
}

// JSONBody 指定一个 JSON 格式的 body
//
// NOTE: 此函并不会设置 content-type 报头。
func (req *Request) JSONBody(obj interface{}) *Request {
	return req.BodyFunc(obj, json.Marshal)
}

// XMLBody 指定一个 XML 格式的 body
//
// NOTE: 此函并不会设置 content-type 报头。
func (req *Request) XMLBody(obj interface{}) *Request {
	return req.BodyFunc(obj, xml.Marshal)
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
