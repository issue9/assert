// SPDX-License-Identifier: MIT

package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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
func (req *Request) Do() *Response {
	req.a.TB().Helper()

	r, err := http.NewRequest(req.method, req.buildPath(), req.body)
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

// Resp 返回 http.Response 实例
//
// NOTE: http.Response.Body 内容已经被读取且关闭。
func (resp *Response) Resp() *http.Response { return resp.resp }

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

// Body 断言内容与 val 相同
func (resp *Response) Body(val []byte, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(bytes.Equal(resp.body, val), msg, []interface{}{"内容并不相同：\nv1=%s\nv2=%s", string(resp.body), string(val)})
}

// StringBody 断言内容与 val 相同
func (resp *Response) StringBody(val string, msg ...interface{}) *Response {
	resp.a.TB().Helper()
	b := string(resp.body)
	return resp.assert(b == val, msg, []interface{}{"内容并不相同：\nv1=%s\nv2=%s", b, val})
}

// BodyNotEmpty 报文内容是否不为空
func (resp *Response) BodyNotEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(len(resp.body) > 0, msg, []interface{}{"内容为空"})
}

// BodyEmpty 报文内容是否为空
func (resp *Response) BodyEmpty(msg ...interface{}) *Response {
	resp.a.TB().Helper()
	return resp.assert(len(resp.body) == 0, msg, []interface{}{"内容并不为空：%s", string(resp.body)})
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
