// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rest

import (
	"net/http"

	"github.com/issue9/assert"
)

// Response 测试请求的返回结构
type Response struct {
	resp      *http.Response
	assertion *assert.Assertion
}

// Success 当前请求是否被成功处理。状态码在 100-399 之间均算成功
func (resp *Response) Success(msg ...interface{}) *Response {
	resp.assertion.True(resp.resp.StatusCode >= 100 && resp.resp.StatusCode < 400)

	return resp
}

// Fail 当前请求是否出错，闫状态码大于 399 均算出错。
func (resp *Response) Fail(msg ...interface{}) *Response {
	resp.assertion.True(resp.resp.StatusCode >= 400)

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
	resp.assertion.Equal(resp.resp.StatusCode, resp.resp.Header.Get(key), msg...)

	return resp
}
