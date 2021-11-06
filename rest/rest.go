// SPDX-License-Identifier: MIT

// Package rest 简单的 API 测试库
package rest

import (
	"net/http"

	"github.com/issue9/assert"
)

// BuildHandler 生成用于测试的 http.Handler 对象
//
// 仅是简单地按以下步骤输出内容：
//  - 输出状态码 code；
//  - 输出报头 headers，以 Add 方式，而不是 set，不会覆盖原来的数据；
//  - 输出 body，如果为空字符串，则不会输出；
func BuildHandler(a *assert.Assertion, code int, body string, headers map[string]string) http.Handler {
	return http.HandlerFunc(BuildHandlerFunc(a, code, body, headers))
}

func BuildHandlerFunc(a *assert.Assertion, code int, body string, headers map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		for k, v := range headers {
			w.Header().Add(k, v)
		}

		if body != "" {
			_, err := w.Write([]byte(body))
			a.NotError(err)
		}
	}
}
