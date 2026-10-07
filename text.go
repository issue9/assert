// SPDX-FileCopyrightText: 2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"fmt"
	"reflect"
)

// EncodingEqual 断言两个编码相同
//
// v1 和 v2 是编码的内容，通过 u 解码为对象 T，如果解码后的两个对象相等则表示两者相等。
// T 的类型不能是指针，否则可能会 panic。
func (a *Assertion) EncodingEqual[T any](v1, v2 []byte, u func(data []byte, v any) error, msg ...any) *Assertion {
	a.TB().Helper()

	v11, v22, err := unmarshal[T](v1, v2, u)
	if err != nil {
		return a.Assert(false, NewFailure("EncodingEqual", msg, map[string]any{"err": err}))
	}

	return a.Assert(isEqual(v11, v22), NewFailure("EncodingEqual", msg, map[string]any{"v1": v1, "v2": v2}))
}

// EncodingNotEqual 断言两个编码不相同
//
// 参数信息可参考 [Assertion.EncodingEqual]。
func (a *Assertion) EncodingNotEqual[T any](v1, v2 []byte, u func(data []byte, v any) error, msg ...any) *Assertion {
	a.TB().Helper()

	v11, v22, err := unmarshal[T](v1, v2, u)
	if err != nil {
		return a.Assert(false, NewFailure("EncodingNotEqual", msg, map[string]any{"err": err}))
	}

	return a.Assert(!isEqual(v11, v22), NewFailure("EncodingNotEqual", msg, map[string]any{"v1": v1, "v2": v2}))
}

func unmarshal[T any](v1, v2 []byte, u func([]byte, any) error) (v111, v222 *T, err error) {
	var v11 T
	k := reflect.TypeFor[T]().Kind()
	if k == reflect.Pointer || k == reflect.Func {
		return nil, nil, fmt.Errorf("类型 T 的 kind %s 无效", k)
	}

	if err := u(v1, &v11); err != nil {
		return nil, nil, err
	}

	var v22 T
	if err := u(v2, &v22); err != nil {
		return nil, nil, err
	}

	return &v11, &v22, nil
}
