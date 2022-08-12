// SPDX-License-Identifier: MIT

// Package rest 简单的 API 测试库
package rest

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/issue9/assert/v3"
)

// BuildHandler 生成用于测试的 http.Handler 对象
//
// 仅是简单地按以下步骤输出内容：
//   - 输出状态码 code；
//   - 输出报头 headers，以 Add 方式，而不是 set，不会覆盖原来的数据；
//   - 输出 body，如果为空字符串，则不会输出；
func BuildHandler(a *assert.Assertion, code int, body string, headers map[string]string) http.Handler {
	return http.HandlerFunc(BuildHandlerFunc(a, code, body, headers))
}

func BuildHandlerFunc(a *assert.Assertion, code int, body string, headers map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a.TB().Helper()

		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.WriteHeader(code)

		if body != "" {
			_, err := w.Write([]byte(body))
			a.NotError(err)
		}
	}
}

func (srv *Server) RawHTTP(req, resp string) *Server {
	srv.Assertion().TB().Helper()
	RawHTTP(srv.Assertion(), srv.client, req, resp)
	return srv
}

// RawHTTP 通过原始数据进行比较请求和返回数据是符合要求
//
// reqRaw 表示原始的请求数据；
// respRaw 表示返回之后的原始数据；
//
// NOTE: 仅判断状态码、报头和实际内容是否相同，而不是直接比较两个 http.Response 的值。
func RawHTTP(a *assert.Assertion, client *http.Client, reqRaw, respRaw string) {
	if client == nil {
		client = &http.Client{}
	}
	a.TB().Helper()

	r, resp := readRaw(a, reqRaw, respRaw)
	if r == nil {
		return
	}

	ret, err := client.Do(r)
	a.NotError(err).NotNil(ret)

	compare(a, resp, ret.StatusCode, ret.Header, ret.Body)
	a.NotError(ret.Body.Close())
}

// RawHandler 通过原始数据进行比较请求和返回数据是符合要求
//
// 功能上与 RawHTTP 相似，处理方式从 http.Client 变成了 http.Handler。
func RawHandler(a *assert.Assertion, h http.Handler, reqRaw, respRaw string) {
	if h == nil {
		panic("h 不能为空")
	}
	a.TB().Helper()

	r, resp := readRaw(a, reqRaw, respRaw)
	if r == nil {
		return
	}

	ret := httptest.NewRecorder()
	h.ServeHTTP(ret, r)

	compare(a, resp, ret.Code, ret.Header(), ret.Body)
}

func readRaw(a *assert.Assertion, reqRaw, respRaw string) (*http.Request, *http.Response) {
	a.TB().Helper()

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewBufferString(respRaw)), nil)
	a.NotError(err).NotNil(resp)

	r, err := http.ReadRequest(bufio.NewReader(bytes.NewBufferString(reqRaw)))
	a.NotError(err).NotNil(r)
	u, err := url.Parse(r.Host + r.URL.String())
	a.NotError(err).NotNil(u)
	r.RequestURI = ""
	r.URL = u

	return r, resp
}

func compare(a *assert.Assertion, resp *http.Response, status int, header http.Header, body io.Reader) {
	a.Equal(resp.StatusCode, status, "compare 断言失败，状态码的期望值 %d 与实际值 %d 不同", resp.StatusCode, status)

	for k := range resp.Header {
		respV := resp.Header.Get(k)
		retV := header.Get(k)
		a.Equal(respV, retV, "compare 断言失败，报头 %s 的期望值 %s 与实际值 %s 不相同", k, respV, retV)
	}

	retB, err := io.ReadAll(body)
	a.NotError(err).NotNil(retB)
	respB, err := io.ReadAll(resp.Body)
	a.NotError(err).NotNil(respB)
	retB = bytes.TrimSpace(retB)
	respB = bytes.TrimSpace(respB)
	a.Equal(respB, retB, "compare 断言失败，内容的期望值与实际值不相同\n%s\n\n%s\n", respB, retB)
}
