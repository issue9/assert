// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"fmt"
	"reflect"
)

// Length 断言长度是否为指定的值
//
// 对于 v 的要求可参考 [Length]
func (a *Assertion) Length(v any, l int, msg ...any) *Assertion {
	a.TB().Helper()

	rl, err := Length(v)
	if err != nil {
		a.Assert(false, NewFailure("Length", msg, map[string]any{"err": err}))
	}
	return a.Assert(rl == l, NewFailure("Length", msg, map[string]any{"l1": rl, "l2": l}))
}

// NotLength 断言长度不是指定的值
//
// 对于 v 的要求可参考 [Length]
func (a *Assertion) NotLength(v any, l int, msg ...any) *Assertion {
	a.TB().Helper()

	rl, err := Length(v)
	if err != nil {
		a.Assert(false, NewFailure("NotLength", msg, map[string]any{"err": err}))
	}
	return a.Assert(rl != l, NewFailure("NotLength", msg, map[string]any{"l": rl}))
}

type RealNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func (a *Assertion) Greater[T RealNumber](v, val T, msg ...any) *Assertion {
	return a.Assert(v > val, NewFailure("Greater", msg, nil))
}

func (a *Assertion) Less[T RealNumber](v, val T, msg ...any) *Assertion {
	return a.Assert(v < val, NewFailure("Less", msg, nil))
}

func (a *Assertion) GreaterEqual[T RealNumber](v, val T, msg ...any) *Assertion {
	return a.Assert(v >= val, NewFailure("GreaterEqual", msg, nil))
}

func (a *Assertion) LessEqual[T RealNumber](v, val T, msg ...any) *Assertion {
	return a.Assert(v <= val, NewFailure("LessEqual", msg, nil))
}

// Positive 断言 v 为大于 0 的数
func (a *Assertion) Positive[T RealNumber](v T, msg ...any) *Assertion {
	return a.Assert(v > 0, NewFailure("Positive", msg, nil))
}

// Negative 断言 v 为小于 0 的数
func (a *Assertion) Negative[T RealNumber](v T, msg ...any) *Assertion {
	return a.Assert(v < 0, NewFailure("Negative", msg, nil))
}

// Between 断言 v 是否存在于 (min,max) 之间
func (a *Assertion) Between[T RealNumber](v, min, max T, msg ...any) *Assertion {
	return a.Assert(v > min && v < max, NewFailure("Between", msg, nil))
}

// BetweenEqual 断言 v 是否存在于 [min,max] 之间
func (a *Assertion) BetweenEqual[T RealNumber](v, min, max T, msg ...any) *Assertion {
	return a.Assert(v >= min && v <= max, NewFailure("BetweenEqual", msg, nil))
}

// BetweenEqualMin 断言 v 是否存在于 [min,max) 之间
func (a *Assertion) BetweenEqualMin[T RealNumber](v, min, max T, msg ...any) *Assertion {
	return a.Assert(v >= min && v < max, NewFailure("BetweenEqualMin", msg, nil))
}

// BetweenEqualMax 断言 v 是否存在于 (min,max] 之间
func (a *Assertion) BetweenEqualMax[T RealNumber](v, min, max T, msg ...any) *Assertion {
	return a.Assert(v > min && v <= max, NewFailure("BetweenEqualMax", msg, nil))
}

// Length 获取对象的长度
//
// v 只能是以下类型或是以下类型的指针，否则将返回错误：
//   - map
//   - string
//   - slice
//   - array
//   - chan
func Length(v any) (int, error) {
	r := reflect.ValueOf(v)
	for r.Kind() == reflect.Pointer {
		r = r.Elem()
	}

	if v == nil {
		return 0, nil
	}

	switch r.Kind() {
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map, reflect.Chan:
		return r.Len(), nil
	}
	return 0, fmt.Errorf("无法获取 %s 类型的长度信息", r.Kind())
}
