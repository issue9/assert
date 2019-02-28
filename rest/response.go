// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/issue9/assert"
)

// Response 测试请求的返回结构
type Response struct {
	resp      *http.Response
	assertion *assert.Assertion
	body      []byte
}

// Do 执行请求操作
func (req *Request) Do() *Response {
	r, err := http.NewRequest(req.method, req.prefix+req.buildPath(), req.body)
	req.assertion.NotError(err).NotNil(r)

	for k, v := range req.headers {
		r.Header.Add(k, v)
	}

	resp, err := req.client.Do(r)
	req.assertion.NotError(err).NotNil(resp)

	bs, err := ioutil.ReadAll(resp.Body)
	req.assertion.NotError(err).NotNil(bs)
	req.assertion.NotError(resp.Body.Close())
	return &Response{
		assertion: req.assertion,
		resp:      resp,
		body:      bs,
	}
}

// Success 当前请求是否被成功处理。状态码在 100-399 之间均算成功
func (resp *Response) Success(msg ...interface{}) *Response {
	resp.assertion.True(resp.resp.StatusCode >= 100 && resp.resp.StatusCode < 400, msg...)

	return resp
}

// Fail 当前请求是否出错，闫状态码大于 399 均算出错。
func (resp *Response) Fail(msg ...interface{}) *Response {
	resp.assertion.True(resp.resp.StatusCode >= 400, msg...)

	return resp
}

// Status 判断状态码是否与 status 相等，若不相等，则返回 msg 指定的消息
//
// msg 可以为空，会返回一个默认的错误提示信息
func (resp *Response) Status(status int, msg ...interface{}) *Response {
	resp.assertion.Equal(resp.resp.StatusCode, status, msg...)

	return resp
}

// NotStatus 判断状态码是否与 status 不相等，若相等，则返回 msg 指定的消息
//
// msg 可以为空，会返回一个默认的错误提示信息
func (resp *Response) NotStatus(status int, msg ...interface{}) *Response {
	resp.assertion.NotEqual(resp.resp.StatusCode, status, msg...)

	return resp
}

// Header 判断指定的报头是否与 val 相同
//
// msg 可以为空，会返回一个默认的错误提示信息
func (resp *Response) Header(key string, val string, msg ...interface{}) *Response {
	resp.assertion.Equal(resp.resp.Header.Get(key), val, msg...)

	return resp
}

// NotHeader 指定的报头必定不与 val 相同。
func (resp *Response) NotHeader(key string, val string, msg ...interface{}) *Response {
	resp.assertion.NotEqual(resp.resp.Header.Get(key), val, msg...)
	return resp
}

// Body 报文内容是否与 val 相等
func (resp *Response) Body(val []byte, msg ...interface{}) *Response {
	resp.assertion.Equal(resp.body, val, msg...)
	return resp
}

// StringBody 报文内容是否与 val 相等
func (resp *Response) StringBody(val string, msg ...interface{}) *Response {
	resp.assertion.Equal(string(resp.body), val, msg...)
	return resp
}

// BodyNotNil 报文内容是否不为 nil
func (resp *Response) BodyNotNil(msg ...interface{}) *Response {
	resp.assertion.NotNil(resp.body, msg...)
	return resp
}

// BodyNil 报文内容是否为 nil
func (resp *Response) BodyNil(msg ...interface{}) *Response {
	resp.assertion.Nil(resp.body, msg...)
	return resp
}

// BodyNotEmpty 报文内容是否不为空
func (resp *Response) BodyNotEmpty(msg ...interface{}) *Response {
	resp.assertion.NotEmpty(resp.body, msg...)
	return resp
}

// BodyEmpty 报文内容是否为空
func (resp *Response) BodyEmpty(msg ...interface{}) *Response {
	resp.assertion.Empty(resp.body, msg...)
	return resp
}

// JSONBody 将 val 转换成 JSON 对象，并与 body 作对比
func (resp *Response) JSONBody(val interface{}) *Response {
	j, err := json.Marshal(val)
	resp.assertion.NotError(err).NotNil(j)

	return resp.Body(j, val)
}

// ReadBody 读取 Body 的内容到 w
func (resp *Response) ReadBody(w io.Writer) *Response {
	n, err := w.Write(resp.body)
	resp.assertion.NotError(err).Equal(n, len(resp.body))
	return resp
}
