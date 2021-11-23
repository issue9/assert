// SPDX-License-Identifier: MIT

package rest

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/issue9/assert/v2"
)

// Response 测试请求的返回结构
type Response struct {
	resp *http.Response
	a    *assert.Assertion
	body []byte
}

// Do 执行请求操作
func (req *Request) Do() *Response {
	req.a.TB().Helper()

	r, err := http.NewRequest(req.method, req.prefix+req.buildPath(), req.body)
	req.a.NotError(err).NotNil(r)

	for k, v := range req.headers {
		r.Header.Add(k, v)
	}

	resp, err := req.client.Do(r)
	req.a.NotError(err).NotNil(resp)

	bs, err := ioutil.ReadAll(resp.Body)
	req.a.NotError(err)
	req.a.NotError(resp.Body.Close())

	return &Response{
		a:    req.a,
		resp: resp,
		body: bs,
	}
}

func (resp *Response) assert(expr bool, msg1, msg2 []interface{}) *Response {
	resp.a.TB().Helper()
	resp.a.Assert(expr, msg1, msg2)
	return resp
}

// Success 状态码是否在 100-399 之间
func (resp *Response) Success(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	succ := resp.resp.StatusCode >= 100 && resp.resp.StatusCode < 400
	return resp.assert(succ, msg, []interface{}{"当前状态码为 %d", resp.resp.StatusCode})
}

// Fail 状态码是否大于 399
func (resp *Response) Fail(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	fail := resp.resp.StatusCode >= 400
	return resp.assert(fail, msg, []interface{}{"当前状态为 %d", resp.resp.StatusCode})
}

// Status 判断状态码是否与 status 相等
func (resp *Response) Status(status int, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	eq := resp.resp.StatusCode == status
	return resp.assert(eq, msg, []interface{}{"实际状态为 %d，与期望值 %d 不符合", resp.resp.StatusCode, status})
}

// NotStatus 判断状态码是否与 status 不相等
func (resp *Response) NotStatus(status int, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	neq := resp.resp.StatusCode != status
	return resp.assert(neq, msg, []interface{}{"状态码 %d 与期望值是相同的", resp.resp.StatusCode})
}

// Header 判断指定的报头是否与 val 相同
//
// msg 可以为空，会返回一个默认的错误提示信息
func (resp *Response) Header(key string, val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	h := resp.resp.Header.Get(key)
	return resp.assert(h == val, msg, []interface{}{"报头 %s 的值 %s 与期望值 %s 不同", key, h, val})
}

// NotHeader 指定的报头必定不与 val 相同。
func (resp *Response) NotHeader(key string, val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	h := resp.resp.Header.Get(key)
	return resp.assert(h != val, msg, []interface{}{"报头 %s 与期望值 %s 相等", key, h})
}

// Body 报文内容是否与 val 相等
func (resp *Response) Body(val []byte, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.Equal(body, val, msg...)
	})
}

// StringBody 报文内容是否与 val 相等
func (resp *Response) StringBody(val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.Equal(string(body), val, msg...)
	})
}

// BodyNotNil 报文内容是否不为 nil
func (resp *Response) BodyNotNil(msg ...interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.NotNil(body, msg...)
	})
}

// BodyNil 报文内容是否为 nil
func (resp *Response) BodyNil(msg ...interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()
		a.Nil(body, msg...)
	})
}

// BodyNotEmpty 报文内容是否不为空
func (resp *Response) BodyNotEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()
		a.NotEmpty(body, msg...)
	})
}

// BodyEmpty 报文内容是否为空
func (resp *Response) BodyEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()
		a.Empty(body, msg...)
	})
}

// JSONBody 将 val 转换成 JSON 对象，并与 body 作对比
func (resp *Response) JSONBody(val interface{}) *Response {
	resp.a.TB().Helper()

	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()
		j, err := json.Marshal(val)
		a.NotError(err).NotNil(j).Equal(body, j)
	})
}

// BodyFunc 指定对 body 内容的断言方式
func (resp *Response) BodyFunc(f func(a *assert.Assertion, body []byte)) *Response {
	resp.a.TB().Helper()
	f(resp.a, resp.body)
	return resp
}

// ReadBody 读取 Body 的内容到 w
func (resp *Response) ReadBody(w io.Writer) *Response {
	n, err := w.Write(resp.body)
	resp.a.NotError(err).Equal(n, len(resp.body))
	return resp
}
