// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"fmt"
	"reflect"
)

// Length 断言长度是否为指定的值
//
// v 可以是以下类型：
//   - map
//   - string
//   - slice
//   - array
func (a *Assertion) Length(v interface{}, l int, msg ...interface{}) *Assertion {
	a.TB().Helper()

	rl, err := getLen(v)
	if err != "" {
		a.Assert(false, NewFailure("Length", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(rl == l, NewFailure("Length", msg, map[string]interface{}{"l1": rl, "l2": l}))
}

// NotLength 断言长度不是指定的值
//
// v 可以是以下类型：
//   - map
//   - string
//   - slice
//   - array
func (a *Assertion) NotLength(v interface{}, l int, msg ...interface{}) *Assertion {
	a.TB().Helper()

	rl, err := getLen(v)
	if err != "" {
		a.Assert(false, NewFailure("NotLength", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(rl != l, NewFailure("NotLength", msg, map[string]interface{}{"l": rl}))
}

func (a *Assertion) Greater(v interface{}, val float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("Greater", msg, nil))
	}
	return a.Assert(vv > val, NewFailure("Greater", msg, nil))
}

func (a *Assertion) Less(v interface{}, val float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("Less", msg, nil))
	}
	return a.Assert(vv < val, NewFailure("Less", msg, nil))
}

func (a *Assertion) GreaterEqual(v interface{}, val float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("GreaterEqual", msg, nil))
	}
	return a.Assert(vv >= val, NewFailure("GreaterEqual", msg, nil))
}

func (a *Assertion) LessEqual(v interface{}, val float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("LessEqual", msg, nil))
	}
	return a.Assert(vv <= val, NewFailure("LessEqual", msg, nil))
}

// Positive 断言 v 为正数
//
// NOTE: 不包含 0
func (a *Assertion) Positive(v interface{}, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("Positive", msg, nil))
	}
	return a.Assert(vv > 0, NewFailure("Positive", msg, nil))
}

// Negative 断言 v 为负数
//
// NOTE: 不包含 0
func (a *Assertion) Negative(v interface{}, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("Negative", msg, nil))
	}
	return a.Assert(vv < 0, NewFailure("Negative", msg, nil))
}

// Between 断言 v 是否存在于 (min,max) 之间
func (a *Assertion) Between(v interface{}, min, max float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("Between", msg, nil))
	}

	return a.Assert(vv > min && vv < max, NewFailure("Between", msg, nil))
}

// BetweenEqual 断言 v 是否存在于 [min,max] 之间
func (a *Assertion) BetweenEqual(v interface{}, min, max float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("BetweenEqual", msg, nil))
	}

	return a.Assert(vv >= min && vv <= max, NewFailure("BetweenEqual", msg, nil))
}

// BetweenEqualMin 断言 v 是否存在于 [min,max) 之间
func (a *Assertion) BetweenEqualMin(v interface{}, min, max float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("BetweenEqualMin", msg, nil))
	}

	return a.Assert(vv >= min && vv < max, NewFailure("BetweenEqualMin", msg, nil))
}

// BetweenEqualMax 断言 v 是否存在于 (min,max] 之间
func (a *Assertion) BetweenEqualMax(v interface{}, min, max float64, msg ...interface{}) *Assertion {
	vv, ok := getNumber(v)
	if !ok {
		return a.Assert(false, NewFailure("BetweenEqualMax", msg, nil))
	}

	return a.Assert(vv > min && vv <= max, NewFailure("BetweenEqualMax", msg, nil))
}

// bool 表示是否成功转换
func getNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return float64(val), true
	}

	return 0, false
}

func getLen(v interface{}) (l int, msg string) {
	r := reflect.ValueOf(v)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	if v == nil {
		return 0, ""
	}

	switch r.Kind() {
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map, reflect.Chan:
		return r.Len(), ""
	}
	return 0, fmt.Sprintf("无法获取 %s 类型的长度信息", r.Kind())
}
