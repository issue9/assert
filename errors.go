// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"fmt"
	"strings"
)

// Error 断言有错误发生
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
//
// [Assertion.NotNil] 的特化版本，限定了类型为 error。
func (a *Assertion) Error(expr error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isNil(expr), NewFailure("Error", msg, map[string]interface{}{"v": expr}))
}

// ErrorString 断言有错误发生且错误信息中包含指定的字符串 str
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
func (a *Assertion) ErrorString(expr error, str string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if isNil(expr) { // 空值，必定没有错误
		return a.Assert(false, NewFailure("ErrorString", msg, map[string]interface{}{"v": expr}))
	}
	return a.Assert(strings.Contains(expr.Error(), str), NewFailure("ErrorString", msg, map[string]interface{}{"v": expr}))
}

// ErrorIs 断言 expr 为 target 类型
//
// 相当于 a.True(errors.Is(expr, target))
func (a *Assertion) ErrorIs(expr, target error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(errors.Is(expr, target), NewFailure("ErrorIs", msg, map[string]interface{}{"err": expr}))
}

// NotError 断言没有错误
//
// [Assertion.Nil] 的特化版本，限定了类型为 error。
func (a *Assertion) NotError(expr error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isNil(expr), NewFailure("NotError", msg, map[string]interface{}{"v": expr}))
}

// Panic 断言函数会发生 panic
func (a *Assertion) Panic(fn func(), msg ...interface{}) *Assertion {
	a.TB().Helper()
	has, _ := hasPanic(fn)
	return a.Assert(has, NewFailure("Panic", msg, nil))
}

// PanicString 断言函数会发生 panic 且 panic 信息中包含指定的字符串内容
func (a *Assertion) PanicString(fn func(), str string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if has, m := hasPanic(fn); has {
		return a.Assert(strings.Contains(fmt.Sprint(m), str), NewFailure("PanicString", msg, map[string]interface{}{"msg": m}))
	}
	return a.Assert(false, NewFailure("PanicString", msg, nil))
}

// PanicType 断言函数会发生 panic 且抛出指定的类型
func (a *Assertion) PanicType(fn func(), typ interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if has, m := hasPanic(fn); has {
		t1, t2 := getType(true, m, typ)
		return a.Assert(t1 == t2, NewFailure("PanicType", msg, map[string]interface{}{"v1": t1, "v2": t2}))
	}
	return a.Assert(false, NewFailure("PanicType", msg, nil))
}

// PanicValue 断言函数会抛出与 v 相同的信息
func (a *Assertion) PanicValue(fn func(), v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if has, m := hasPanic(fn); has {
		return a.Assert(isEqual(m, v), NewFailure("PanicValue", msg, map[string]interface{}{"v": m}))
	}
	return a.Assert(false, NewFailure("PanicType", msg, nil))
}

// NotPanic 断言 fn 不会 panic
func (a *Assertion) NotPanic(fn func(), msg ...interface{}) *Assertion {
	a.TB().Helper()
	has, m := hasPanic(fn)
	return a.Assert(!has, NewFailure("NotPanic", msg, map[string]interface{}{"err": m}))
}

// hasPanic 判断 fn 函数是否会发生 panic
// 若发生了 panic，将把 msg 一起返回。
func hasPanic(fn func()) (has bool, msg interface{}) {
	defer func() {
		if msg = recover(); msg != nil {
			has = true
		}
	}()
	fn()

	return
}
