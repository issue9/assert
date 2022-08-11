// SPDX-License-Identifier: MIT

package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/issue9/assert/v2"
)

// Response 测试请求的返回结构
type Response struct {
	resp *http.Response
	a    *assert.Assertion
	body []byte
}

// Do 执行请求操作
//
// h 默认为空，如果不为空，则表示当前请求忽略 http.Client，而是访问 h.ServeHTTP 的内容。
func (req *Request) Do(h http.Handler) *Response {
	if req.client == nil && h == nil {
		panic("h 不能为空")
	}

	req.a.TB().Helper()

	r := req.Request()
	var err error
	var resp *http.Response
	if h != nil {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		resp = w.Result()
	} else {
		resp, err = req.client.Do(r)
		req.a.NotError(err).NotNil(resp)
	}

	var bs []byte
	if resp.Body != nil {
		bs, err = io.ReadAll(resp.Body)
		if err != io.EOF {
			req.a.NotError(err)
		}
		req.a.NotError(resp.Body.Close())
	}

	return &Response{
		a:    req.a,
		resp: resp,
		body: bs,
	}
}

// Resp 返回 http.Response 实例
//
// NOTE: http.Response.Body 内容已经被读取且关闭。
func (resp *Response) Resp() *http.Response { return resp.resp }

func (resp *Response) assert(expr bool, f *assert.Failure) *Response {
	resp.a.TB().Helper()
	resp.a.Assert(expr, f)
	return resp
}

// Success 状态码是否在 100-399 之间
func (resp *Response) Success(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	succ := resp.resp.StatusCode >= 100 && resp.resp.StatusCode < 400
	return resp.assert(succ, assert.NewFailure("Success", msg, map[string]interface{}{"status": resp.resp.StatusCode}))
}

// Fail 状态码是否大于 399
func (resp *Response) Fail(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	fail := resp.resp.StatusCode >= 400
	return resp.assert(fail, assert.NewFailure("Fail", msg, map[string]interface{}{"status": resp.resp.StatusCode}))
}

// Status 判断状态码是否与 status 相等
func (resp *Response) Status(status int, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	eq := resp.resp.StatusCode == status
	return resp.assert(eq, assert.NewFailure("Status", msg, map[string]interface{}{"status1": resp.resp.StatusCode, "status2": status}))
}

// NotStatus 判断状态码是否与 status 不相等
func (resp *Response) NotStatus(status int, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	neq := resp.resp.StatusCode != status
	return resp.assert(neq, assert.NewFailure("NotStatus", msg, map[string]interface{}{"status": resp.resp.StatusCode}))
}

// Header 判断指定的报头是否与 val 相同
//
// msg 可以为空，会返回一个默认的错误提示信息
func (resp *Response) Header(key string, val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	h := resp.resp.Header.Get(key)
	return resp.assert(h == val, assert.NewFailure("Header", msg, map[string]interface{}{"header": key, "v1": h, "v2": val}))
}

// NotHeader 指定的报头必定不与 val 相同。
func (resp *Response) NotHeader(key string, val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	h := resp.resp.Header.Get(key)
	return resp.assert(h != val, assert.NewFailure("NotHeader", msg, map[string]interface{}{"header": key, "v": h}))
}

// Body 断言内容与 val 相同
func (resp *Response) Body(val []byte, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(bytes.Equal(resp.body, val), assert.NewFailure("Body", msg, map[string]interface{}{"v1": string(resp.body), "v2": string(val)}))
}

// StringBody 断言内容与 val 相同
func (resp *Response) StringBody(val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	b := string(resp.body)
	return resp.assert(b == val, assert.NewFailure("StringBody", msg, map[string]interface{}{"v1": b, "v2": val}))
}

// BodyNotEmpty 报文内容是否不为空
func (resp *Response) BodyNotEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(len(resp.body) > 0, assert.NewFailure("BodyNotEmpty", msg, nil))
}

// BodyEmpty 报文内容是否为空
func (resp *Response) BodyEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(len(resp.body) == 0, assert.NewFailure("BodyEmpty", msg, map[string]interface{}{"v": resp.body}))
}

// JSONBody body 转换成 JSON 对象之后是否等价于 val
func (resp *Response) JSONBody(val interface{}) *Response {
	resp.a.TB().Helper()
	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()

		// NOTE: 应当始终将 body 转换 val 相同的类型，然后再比较对象，
		// 因为 val 转换成字符串，可能因为空格缩进等原因，未必会与 body 是相同的。
		b, err := UnmarshalObject(body, val, json.Unmarshal)
		a.NotError(err).Equal(b, val)
	})
}

// XMLBody body 转换成 XML 对象之后是否等价于 val
func (resp *Response) XMLBody(val interface{}) *Response {
	resp.a.TB().Helper()
	return resp.BodyFunc(func(a *assert.Assertion, body []byte) {
		a.TB().Helper()

		// NOTE: 应当始终将 body 转换 val 相同的类型，然后再比较对象，
		// 因为 val 转换成字符串，可能因为空格缩进等原因，未必会与 body 是相同的。
		b, err := UnmarshalObject(body, val, xml.Unmarshal)
		a.NotError(err).Equal(b, val)
	})
}

// BodyFunc 指定对 body 内容的断言方式
func (resp *Response) BodyFunc(f func(a *assert.Assertion, body []byte)) *Response {
	resp.a.TB().Helper()

	b := make([]byte, len(resp.body))
	copy(b, resp.body)
	f(resp.a, b)

	return resp
}

// UnmarshalObject 将 data 以 u 作为转换方式转换成与 val 相同的类型
//
// 如果 val 是指针，则会转换成其指向的类型，返回的对象是指针类型。
func UnmarshalObject(data []byte, val interface{}, u func([]byte, interface{}) error) (interface{}, error) {
	t := reflect.TypeOf(val)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	bv := reflect.New(t)

	if err := u(data, bv.Interface()); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return bv.Interface(), nil
}
