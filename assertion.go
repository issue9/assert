// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

// Assertion是对testing.T和testing.B进行了简单的封装。
// 可以以对象的方式调用包中的各个断言函数。
type Assertion struct {
	t tester
}

// 返回Assertion对象。
func New(t tester) *Assertion {
	return &Assertion{t: t}
}

// 返回tester接口，该接口包含了testing.T和testing.B的共有接口
func (a *Assertion) T() tester {
	return a.t
}

// 参照assert.True()函数
func (a *Assertion) True(expr bool, msg ...interface{}) *Assertion {
	True(a.t, expr, msg...)
	return a
}

// 参照assert.False()函数
func (a *Assertion) False(expr bool, msg ...interface{}) *Assertion {
	False(a.t, expr, msg...)
	return a
}

// 参照assert.Nil()函数
func (a *Assertion) Nil(expr interface{}, msg ...interface{}) *Assertion {
	Nil(a.t, expr, msg...)
	return a
}

// 参照assert.NotNil()函数
func (a *Assertion) NotNil(expr interface{}, msg ...interface{}) *Assertion {
	NotNil(a.t, expr, msg...)
	return a
}

// 参照assert.Equal()函数
func (a *Assertion) Equal(v1, v2 interface{}, msg ...interface{}) *Assertion {
	Equal(a.t, v1, v2, msg...)
	return a
}

// 参照assert.NotEqual()函数
func (a *Assertion) NotEqual(v1, v2 interface{}, msg ...interface{}) *Assertion {
	NotEqual(a.t, v1, v2, msg...)
	return a
}

// 参照assert.Empty()函数
func (a *Assertion) Empty(expr interface{}, msg ...interface{}) *Assertion {
	Empty(a.t, expr, msg...)
	return a
}

// 参照assert.NotEmpty()函数
func (a *Assertion) NotEmpty(expr interface{}, msg ...interface{}) *Assertion {
	NotEmpty(a.t, expr, msg...)
	return a
}

// 参照assert.Error()函数
func (a *Assertion) Error(expr interface{}, msg ...interface{}) *Assertion {
	Error(a.t, expr, msg...)
	return a
}

// 参照assert.ErrorString()函数
func (a *Assertion) ErrorString(expr interface{}, str string, msg ...interface{}) *Assertion {
	ErrorString(a.t, expr, str, msg...)
	return a
}

// 参照assert.ErrorType()函数
func (a *Assertion) ErrorType(expr interface{}, typ error, msg ...interface{}) *Assertion {
	ErrorType(a.t, expr, typ, msg...)
	return a
}

// 参照assert.NotError()函数
func (a *Assertion) NotError(expr interface{}, msg ...interface{}) *Assertion {
	NotError(a.t, expr, msg...)
	return a
}

// 参照assert.FileExists()函数
func (a *Assertion) FileExists(path string, msg ...interface{}) *Assertion {
	FileExists(a.t, path, msg...)
	return a
}

// 参照assert.FileNotExists()函数
func (a *Assertion) FileNotExists(path string, msg ...interface{}) *Assertion {
	FileNotExists(a.t, path, msg...)
	return a
}

// 参照assert.Panic()函数
func (a *Assertion) Panic(fn func(), msg ...interface{}) *Assertion {
	Panic(a.t, fn, msg...)
	return a
}

// 参照assert.PanicString()函数
func (a *Assertion) PanicString(fn func(), str string, msg ...interface{}) *Assertion {
	PanicString(a.t, fn, str, msg...)
	return a
}

// 参照assert.PanicType()函数
func (a *Assertion) PanicType(fn func(), typ interface{}, msg ...interface{}) *Assertion {
	PanicType(a.t, fn, typ, msg...)
	return a
}

// 参照assert.NotPanic()函数
func (a *Assertion) NotPanic(fn func(), msg ...interface{}) *Assertion {
	NotPanic(a.t, fn, msg...)
	return a
}

// 参照assert.Contains()函数
func (a *Assertion) Contains(container, item interface{}, msg ...interface{}) *Assertion {
	Contains(a.t, container, item, msg...)
	return a
}

// 参照assert.NotContains()函数
func (a *Assertion) NotContains(container, item interface{}, msg ...interface{}) *Assertion {
	NotContains(a.t, container, item, msg...)
	return a
}
