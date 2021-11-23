// SPDX-License-Identifier: MIT

package assert

import "testing"

// Assertion 可以以对象的方式调用包中的各个断言函数
type Assertion struct {
	tb    testing.TB
	suite suite
}

// New 返回 Assertion 对象。
func New(tb testing.TB) *Assertion { return &Assertion{tb: tb} }

// Assert 断言 expr 条件成立
//
// expr 返回结果值为 bool 类型的表达式；
// msg1,msg2 输出的错误信息，之所以提供两组信息，是方便在用户没有提供的情况下，
// 可以使用系统内部提供的信息，优先使用 msg1 中的信息，若不存在，则使用 msg2 的内容。
//
// 直接使用 True 断言效果是一样的，之所以提供该函数，
// 主要供库调用，可以提供一个默认的错误信息。
func (a *Assertion) Assert(expr bool, msg1, msg2 []interface{}) *Assertion {
	Assert(a.TB(), expr, msg1, msg2)
	return a
}

// TB 返回 testing.TB 接口
func (a *Assertion) TB() testing.TB { return a.tb }

// True 参照 assert.True() 函数
func (a *Assertion) True(expr bool, msg ...interface{}) *Assertion {
	True(a.tb, expr, msg...)
	return a
}

// False 参照 assert.False() 函数
func (a *Assertion) False(expr bool, msg ...interface{}) *Assertion {
	False(a.tb, expr, msg...)
	return a
}

// Nil 参照 assert.Nil() 函数
func (a *Assertion) Nil(expr interface{}, msg ...interface{}) *Assertion {
	Nil(a.tb, expr, msg...)
	return a
}

// NotNil 参照 assert.NotNil() 函数
func (a *Assertion) NotNil(expr interface{}, msg ...interface{}) *Assertion {
	NotNil(a.tb, expr, msg...)
	return a
}

// Equal 参照 assert.Equal() 函数
func (a *Assertion) Equal(v1, v2 interface{}, msg ...interface{}) *Assertion {
	Equal(a.tb, v1, v2, msg...)
	return a
}

// NotEqual 参照 assert.NotEqual() 函数
func (a *Assertion) NotEqual(v1, v2 interface{}, msg ...interface{}) *Assertion {
	NotEqual(a.tb, v1, v2, msg...)
	return a
}

// Empty 参照 assert.Empty() 函数
func (a *Assertion) Empty(expr interface{}, msg ...interface{}) *Assertion {
	Empty(a.tb, expr, msg...)
	return a
}

// NotEmpty 参照 assert.NotEmpty() 函数
func (a *Assertion) NotEmpty(expr interface{}, msg ...interface{}) *Assertion {
	NotEmpty(a.tb, expr, msg...)
	return a
}

// Error 参照 assert.Error() 函数
func (a *Assertion) Error(expr interface{}, msg ...interface{}) *Assertion {
	Error(a.tb, expr, msg...)
	return a
}

// ErrorString 参照 assert.ErrorString() 函数
func (a *Assertion) ErrorString(expr interface{}, str string, msg ...interface{}) *Assertion {
	ErrorString(a.tb, expr, str, msg...)
	return a
}

// ErrorType 参照 assert.ErrorType() 函数
func (a *Assertion) ErrorType(expr interface{}, typ error, msg ...interface{}) *Assertion {
	ErrorType(a.tb, expr, typ, msg...)
	return a
}

// NotError 参照 assert.NotError() 函数
func (a *Assertion) NotError(expr interface{}, msg ...interface{}) *Assertion {
	NotError(a.tb, expr, msg...)
	return a
}

// ErrorIs 断言 expr 为 target 类型
//
// 相当于 a.True(errors.Is(expr, target))
func (a *Assertion) ErrorIs(expr interface{}, target error, msg ...interface{}) *Assertion {
	ErrorIs(a.tb, expr, target, msg...)
	return a
}

// FileExists 参照 assert.FileExists() 函数
func (a *Assertion) FileExists(path string, msg ...interface{}) *Assertion {
	FileExists(a.tb, path, msg...)
	return a
}

// FileNotExists 参照 assert.FileNotExists() 函数
func (a *Assertion) FileNotExists(path string, msg ...interface{}) *Assertion {
	FileNotExists(a.tb, path, msg...)
	return a
}

// Panic 参照 assert.Panic() 函数
func (a *Assertion) Panic(fn func(), msg ...interface{}) *Assertion {
	Panic(a.tb, fn, msg...)
	return a
}

// PanicString 参照 assert.PanicString() 函数
func (a *Assertion) PanicString(fn func(), str string, msg ...interface{}) *Assertion {
	PanicString(a.tb, fn, str, msg...)
	return a
}

// PanicType 参照 assert.PanicType() 函数
func (a *Assertion) PanicType(fn func(), typ interface{}, msg ...interface{}) *Assertion {
	PanicType(a.tb, fn, typ, msg...)
	return a
}

// NotPanic 参照 assert.NotPanic() 函数
func (a *Assertion) NotPanic(fn func(), msg ...interface{}) *Assertion {
	NotPanic(a.tb, fn, msg...)
	return a
}

// Contains 参照 assert.Contains() 函数
func (a *Assertion) Contains(container, item interface{}, msg ...interface{}) *Assertion {
	Contains(a.tb, container, item, msg...)
	return a
}

// NotContains 参照 assert.NotContains() 函数
func (a *Assertion) NotContains(container, item interface{}, msg ...interface{}) *Assertion {
	NotContains(a.tb, container, item, msg...)
	return a
}

// Zero 断言是否为零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func (a *Assertion) Zero(v interface{}, msg ...interface{}) *Assertion {
	Zero(a.tb, v, msg...)
	return a
}

// NotZero 断言是否为非零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func (a *Assertion) NotZero(v interface{}, msg ...interface{}) *Assertion {
	NotZero(a.tb, v, msg...)
	return a
}

func (a *Assertion) Length(v interface{}, l int, msg ...interface{}) *Assertion {
	Length(a.tb, v, l, msg)
	return a
}

func (a *Assertion) NotLength(v interface{}, l int, msg ...interface{}) *Assertion {
	NotLength(a.tb, v, l, msg)
	return a
}
