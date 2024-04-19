// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"
)

// Assertion 是对 [testing.TB] 的二次包装
type Assertion struct {
	tb    testing.TB
	print func(...interface{})
}

// New 返回 [Assertion] 对象
//
// fatal 决定在出错时是调用 [testing.TB.Error] 还是 [testing.TB.Fatal]；
func New(tb testing.TB, fatal bool) *Assertion {
	p := tb.Error
	if fatal {
		p = tb.Fatal
	}

	return &Assertion{
		tb:    tb,
		print: p,
	}
}

// NewWithEnv 以指定的环境变量初始化 [Assertion] 对象
//
// env 是以 [testing.TB.Setenv] 的形式调用。
func NewWithEnv(tb testing.TB, fatal bool, env map[string]string) *Assertion {
	for k, v := range env {
		tb.Setenv(k, v)
	}
	return New(tb, fatal)
}

// Assert 断言 expr 条件成立
//
// f 表示在断言失败时输出的信息
//
// 普通用户直接使用 [Assertion.True] 效果是一样的，此函数主要供 [Assertion] 自身调用。
func (a *Assertion) Assert(expr bool, f *Failure) *Assertion {
	if !expr {
		a.TB().Helper()
		a.print(GetFailureSprintFunc()(f))
	}
	failurePool.Put(f)
	return a
}

// TB 返回 [testing.TB] 接口
func (a *Assertion) TB() testing.TB { return a.tb }

// True 断言表达式 expr 为真
//
// args 对应 [fmt.Printf] 函数中的参数，其中 args[0] 对应第一个参数 format，依次类推，
// 其它断言函数的 args 参数，功能与此相同。
func (a *Assertion) True(expr bool, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(expr, NewFailure("True", msg, nil))
}

func (a *Assertion) False(expr bool, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!expr, NewFailure("False", msg, nil))
}

func (a *Assertion) Nil(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isNil(expr), NewFailure("Nil", msg, map[string]interface{}{"v": expr}))
}

func (a *Assertion) NotNil(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isNil(expr), NewFailure("NotNil", msg, map[string]interface{}{"v": expr}))
}

func (a *Assertion) Equal(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isEqual(v1, v2), NewFailure("Equal", msg, map[string]interface{}{"v1": v1, "v2": v2}))
}

func (a *Assertion) NotEqual(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isEqual(v1, v2), NewFailure("NotEqual", msg, map[string]interface{}{"v1": v1, "v2": v2}))
}

func (a *Assertion) Empty(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isEmpty(expr), NewFailure("Empty", msg, map[string]interface{}{"v": expr}))
}

func (a *Assertion) NotEmpty(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isEmpty(expr), NewFailure("NotEmpty", msg, map[string]interface{}{"v": expr}))
}

// Contains 断言 container 包含 item 或是包含 item 中的所有项
//
// 若 container string、[]byte 和 []rune 类型，
// 都将会以字符串的形式判断其是否包含 item。
// 若 container 是个列表(array、slice、map)则判断其元素中是否包含 item 中的
// 的所有项，或是 item 本身就是 container 中的一个元素。
func (a *Assertion) Contains(container, item interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isContains(container, item), NewFailure("Contains", msg, map[string]interface{}{"container": container, "item": item}))
}

// NotContains 断言 container 不包含 item 或是不包含 item 中的所有项
func (a *Assertion) NotContains(container, item interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isContains(container, item), NewFailure("NotContains", msg, map[string]interface{}{"container": container, "item": item}))
}

// Zero 断言是否为零值
//
// 最终调用的是 [reflect.Value.IsZero] 进行判断，如果是指针，则会判断指向的对象。
func (a *Assertion) Zero(v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isZero(v), NewFailure("Zero", msg, map[string]interface{}{"v": v}))
}

// NotZero 断言是否为非零值
//
// 最终调用的是 [reflect.Value.IsZero] 进行判断，如果是指针，则会判断指向的对象。
func (a *Assertion) NotZero(v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isZero(v), NewFailure("NotZero", msg, map[string]interface{}{"v": v}))
}

// TypeEqual 断言两个值的类型是否相同
//
// ptr 如果为 true，则会在对象为指针时，查找其指向的对象。
func (a *Assertion) TypeEqual(ptr bool, v1, v2 interface{}, msg ...interface{}) *Assertion {
	if v1 == v2 {
		return a
	}

	a.TB().Helper()

	t1, t2 := getType(ptr, v1, v2)
	return a.Assert(t1 == t2, NewFailure("TypeEqual", msg, map[string]interface{}{"v1": t1, "v2": t2}))
}

// Same 断言为同一个对象
func (a *Assertion) Same(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(isSame(v1, v2), NewFailure("Same", msg, nil))
}

// NotSame 断言为不是同一个对象
func (a *Assertion) NotSame(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!isSame(v1, v2), NewFailure("NotSame", msg, nil))
}

func isSame(v1, v2 interface{}) bool {
	rv1 := reflect.ValueOf(v1)
	if !canPointer(rv1.Kind()) {
		return false
	}
	rv2 := reflect.ValueOf(v2)
	if !canPointer(rv2.Kind()) {
		return false
	}

	return rv1.Pointer() == rv2.Pointer()
}

func canPointer(k reflect.Kind) bool {
	switch k {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice, reflect.UnsafePointer, reflect.Func:
		return true
	default:
		return false
	}
}

// Match 断言 v 是否匹配正则表达式 reg
func (a *Assertion) Match(reg *regexp.Regexp, v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	switch val := v.(type) {
	case string:
		return a.Assert(reg.MatchString(val), NewFailure("Match", msg, map[string]interface{}{"v": val}))
	case []byte:
		return a.Assert(reg.Match(val), NewFailure("Match", msg, map[string]interface{}{"v": val}))
	default:
		return a.Assert(reg.MatchString(fmt.Sprint(val)), NewFailure("Match", msg, map[string]interface{}{"v": val}))
	}
}

// NotMatch 断言 v 是否不匹配正则表达式 reg
func (a *Assertion) NotMatch(reg *regexp.Regexp, v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	switch val := v.(type) {
	case string:
		return a.Assert(!reg.MatchString(val), NewFailure("NotMatch", msg, map[string]interface{}{"v": val}))
	case []byte:
		return a.Assert(!reg.Match(val), NewFailure("NotMatch", msg, map[string]interface{}{"v": val}))
	default:
		return a.Assert(!reg.MatchString(fmt.Sprint(val)), NewFailure("NotMatch", msg, map[string]interface{}{"v": val}))
	}
}

// When 断言 expr 为 true 且在条件成立时调用 f
//
// 当有一组依赖 expr 的断言时，可以调用此方法。f 的参数 a 即为当前实例。
func (a *Assertion) When(expr bool, f func(a *Assertion), msg ...interface{}) *Assertion {
	if expr {
		f(a)
	} else {
		a.TB().Helper()
		a.Assert(false, NewFailure("When", msg, nil))
	}
	return a
}

// Wait 等待一定时间再执行后续操作
func (a *Assertion) Wait(d time.Duration) *Assertion {
	time.Sleep(d)
	return a
}

// WaitSeconds 等待 s 秒再执行后续操作
func (a *Assertion) WaitSeconds(s int) *Assertion { return a.Wait(time.Duration(s) * time.Second) }

// Go 以 goroutine 方式执行 f
func (a *Assertion) Go(f func(*Assertion)) *Assertion {
	go f(a)
	return a
}
