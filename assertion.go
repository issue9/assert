// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"testing/synctest"
	"time"

	"github.com/issue9/assert/v5/internal"
)

// Assertion 是对 [testing.TB] 的二次包装
type Assertion struct {
	tb  testing.TB
	log func(...any)
	env map[string]string
}

// New 创建 [Assertion] 对象
//
// fatal 决定在出错时是调用 [testing.TB.Error] 还是 [testing.TB.Fatal]；
func New(tb testing.TB, fatal bool) *Assertion {
	return NewWithEnv(tb, fatal, nil)
}

// NewWithEnv 以指定的环境变量初始化 [Assertion] 对象
//
// fatal 决定在出错时是调用 [testing.TB.Error] 还是 [testing.TB.Fatal]；
// env 是以 [testing.TB.Setenv] 的形式调用；
func NewWithEnv(tb testing.TB, fatal bool, env map[string]string) *Assertion {
	p := tb.Error
	if fatal {
		p = tb.Fatal
	}

	return newWithLogEnv(tb, p, env)
}

func newWithLogEnv(tb testing.TB, log func(...any), env map[string]string) *Assertion {
	for k, v := range env {
		tb.Setenv(k, v)
	}

	return &Assertion{
		tb:  tb,
		log: log,
		env: env,
	}
}

// Assert 断言 expr 条件成立
//
// f 表示在断言失败时输出的信息
//
// 普通用户直接使用 [Assertion.True] 效果是一样的，此函数主要供 [Assertion] 自身调用。
func (a *Assertion) Assert(expr bool, f *Failure) *Assertion {
	if !expr {
		a.TB().Helper()

		// 如果这里调用了 Fail，那么不触发 failurePool.Put 回收 f，sync.Pool 不回收不会造成内存泄漏。
		a.log(GetFailureSprintFunc()(f))
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
func (a *Assertion) True(expr bool, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(expr, NewFailure("True", msg, nil))
}

func (a *Assertion) False(expr bool, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!expr, NewFailure("False", msg, nil))
}

func (a *Assertion) Nil(expr any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(internal.IsNil(expr), NewFailure("Nil", msg, map[string]any{"v": expr}))
}

func (a *Assertion) NotNil(expr any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!internal.IsNil(expr), NewFailure("NotNil", msg, map[string]any{"v": expr}))
}

// Equal 判断两值值是否相等
//
// 该方法法需要两值的类型是相等的，比如 int(8) 和 int16(8) 是相等的，
// 甚至 map[string]int{"key":5} 和 map[string]int8{"key":5} 也是相等的。
func (a *Assertion) Equal(v1, v2 any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(internal.IsEqual(v1, v2), NewFailure("Equal", msg, map[string]any{"v1": v1, "v2": v2}))
}

// StrictEqual 从类型到值都相等
func (a *Assertion) StrictEqual[T comparable](v1, v2 T, msg ...any) *Assertion {
	return a.Assert(v1 == v2, NewFailure("StrictEqual", msg, map[string]any{"v1": v1, "v2": v2}))
}

func (a *Assertion) NotEqual(v1, v2 any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!internal.IsEqual(v1, v2), NewFailure("NotEqual", msg, map[string]any{"v1": v1, "v2": v2}))
}

// Empty 判断对象是否为空
//
// 与 [Assertion.Zero] 相比，包含了对容器对象的长度为 0 的判断。
func (a *Assertion) Empty(expr any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(internal.IsEmpty(expr), NewFailure("Empty", msg, map[string]any{"v": expr}))
}

func (a *Assertion) NotEmpty(expr any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!internal.IsEmpty(expr), NewFailure("NotEmpty", msg, map[string]any{"v": expr}))
}

// Contains 断言 container 包含 item 或是包含 item 中的所有项
//
// 若 container 是 string、[]byte 和 []rune 类型，
// 都将会以字符串的形式判断其是否包含 item。
// 若 container 是个列表(array、slice、map)则判断其元素中是否包含 item 中的
// 的所有项，或是 item 本身就是 container 中的一个元素。
func (a *Assertion) Contains(container, item any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(internal.IsContains(container, item), NewFailure("Contains", msg, map[string]any{"container": container, "item": item}))
}

// NotContains 断言 container 不包含 item 或是不包含 item 中的所有项
func (a *Assertion) NotContains(container, item any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!internal.IsContains(container, item), NewFailure("NotContains", msg, map[string]any{"container": container, "item": item}))
}

// Zero 断言是否为零值
//
// 最终调用的是 [reflect.Value.IsZero] 进行判断，如果是指针，则会判断指向的对象。
func (a *Assertion) Zero(v any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(internal.IsZero(v), NewFailure("Zero", msg, map[string]any{"v": v}))
}

// NotZero 断言是否为非零值
//
// 最终调用的是 [reflect.Value.IsZero] 进行判断，如果是指针，则会判断指向的对象。
func (a *Assertion) NotZero(v any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!internal.IsZero(v), NewFailure("NotZero", msg, map[string]any{"v": v}))
}

// TypeEqual 断言两个值的类型是否相同
//
// ptr 如果为 true，则会在对象为指针时，查找其指向的对象。
func (a *Assertion) TypeEqual(ptr bool, v1, v2 any, msg ...any) *Assertion {
	if v1 == v2 {
		return a
	}

	a.TB().Helper()

	t1, t2 := internal.GetType(ptr, v1, v2)
	return a.Assert(t1 == t2, NewFailure("TypeEqual", msg, map[string]any{"v1": t1, "v2": t2}))
}

// Same 断言为同一个对象
func (a *Assertion) Same[T any](v1, v2 T, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(isSame(v1, v2), NewFailure("Same", msg, nil))
}

// NotSame 断言为不是同一个对象
func (a *Assertion) NotSame(v1, v2 any, msg ...any) *Assertion {
	a.TB().Helper()
	return a.Assert(!isSame(v1, v2), NewFailure("NotSame", msg, nil))
}

func isSame(v1, v2 any) bool {
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
	case reflect.Pointer, reflect.Map, reflect.Chan, reflect.Slice, reflect.UnsafePointer, reflect.Func:
		return true
	default:
		return false
	}
}

// Match 断言 v 是否匹配正则表达式 reg
func (a *Assertion) Match(reg *regexp.Regexp, v any, msg ...any) *Assertion {
	a.TB().Helper()
	switch val := v.(type) {
	case string:
		return a.Assert(reg.MatchString(val), NewFailure("Match", msg, map[string]any{"v": val}))
	case []byte:
		return a.Assert(reg.Match(val), NewFailure("Match", msg, map[string]any{"v": val}))
	default:
		return a.Assert(reg.MatchString(fmt.Sprint(val)), NewFailure("Match", msg, map[string]any{"v": val}))
	}
}

// NotMatch 断言 v 是否不匹配正则表达式 reg
func (a *Assertion) NotMatch(reg *regexp.Regexp, v any, msg ...any) *Assertion {
	a.TB().Helper()
	switch val := v.(type) {
	case string:
		return a.Assert(!reg.MatchString(val), NewFailure("NotMatch", msg, map[string]any{"v": val}))
	case []byte:
		return a.Assert(!reg.Match(val), NewFailure("NotMatch", msg, map[string]any{"v": val}))
	default:
		return a.Assert(!reg.MatchString(fmt.Sprint(val)), NewFailure("NotMatch", msg, map[string]any{"v": val}))
	}
}

// When 断言 expr 为 true 且在条件成立时调用 f
//
// 当有一组依赖 expr 的断言时，可以调用此方法。f 的参数 a 即为当前实例。
func (a *Assertion) When(expr bool, f func(a *Assertion), msg ...any) *Assertion {
	if expr {
		f(a)
	}
	return a
}

// SyncTest 调用 [synctest.Test] 的测试用例
//
// NOTE: 此方法要求 [Assertion.TB] 的类型必须为 [testing.T]。
func (a *Assertion) SyncTest(f func(a *Assertion, wait func())) *Assertion {
	t, ok := a.TB().(*testing.T)
	if !ok {
		panic("SyncTest 只能应用在 testing.T 上")
	}

	synctest.Test(t, func(t *testing.T) { f(newWithLogEnv(t, a.log, a.env), synctest.Wait) })

	return a
}

// Eventually 断言在 timeout 时间之内 f 会返回 true
func (a *Assertion) Eventually(f func() bool, timeout time.Duration, msg ...any) *Assertion {
	a.TB().Helper()
	return a.try(false, f, timeout, msg...)
}

// Never 断言在 timeout 时间之内 f 始终返回 false
func (a *Assertion) Never(f func() bool, timeout time.Duration, msg ...any) *Assertion {
	a.TB().Helper()
	return a.try(true, f, timeout, msg...)
}

func (a *Assertion) try(never bool, f func() bool, timeout time.Duration, msg ...any) *Assertion {
	a.TB().Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // cancel 多次调用不影响

	ret := make(chan bool, 1) // 保存异步调用的返回值
	ff := func() { ret <- f() }

	go ff()

LOOP:
	for {
		select {
		case <-ctx.Done():
			if never {
				// Never 模式下，如果是 cancel 表示断言失败
				a.Assert(errors.Is(ctx.Err(), context.DeadlineExceeded), NewFailure("Never", msg, nil))
			} else {
				// Eventually 模式下，如果是 Deadline 表示断言失败
				a.Assert(errors.Is(ctx.Err(), context.Canceled), NewFailure("Eventually", msg, nil))
			}
			break LOOP
		case v := <-ret:
			if v { // 无论是 never 还是 eventually 都是返回 true 退出
				cancel()
			} else {
				go ff() // 收到 ff 的返回值，才进行下一次调用。
			}
		}
	}

	return a
}
