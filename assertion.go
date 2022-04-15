// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Assertion 可以以对象的方式调用包中的各个断言函数
type Assertion struct {
	tb testing.TB

	fatal  bool
	print  func(...interface{})
	printf func(string, ...interface{})

	suite suite
}

// New 返回 Assertion 对象
//
// fatal 决定在出错时是调用 tb.Error 还是 tb.Fatal；
func New(tb testing.TB, fatal bool) *Assertion {
	p := tb.Error
	pf := tb.Errorf
	if fatal {
		p = tb.Fatal
		pf = tb.Fatalf
	}

	return &Assertion{
		tb: tb,

		fatal:  fatal,
		print:  p,
		printf: pf,
	}
}

// Assert 断言 expr 条件成立
//
// msg1,msg2 输出的错误信息，优先使用 msg1 中的信息，若不存在，则使用 msg2 的内容。
// msg1 和 msg2 格式完全相同，根据其每一个元素是否为 string 决定是调用 Error 还是 Errorf。
//
// 普通用户直接使用 True 效果是一样的，
// 之所以提供该函数，主要供库调用，可以提供一个默认的错误信息。
func (a *Assertion) Assert(expr bool, msg1, msg2 []interface{}) *Assertion {
	if expr {
		return a
	}

	a.TB().Helper()

	if len(msg1) == 0 {
		msg1 = msg2
	}
	if len(msg1) == 0 {
		panic("未提供任何错误信息")
	}

	if len(msg1) == 1 {
		a.print(msg1...)
		return a
	}

	if format, ok := msg1[0].(string); ok {
		a.printf(format, msg1[1:]...)
	} else {
		a.print(msg1...)
	}

	return a
}

// TB 返回 testing.TB 接口
func (a *Assertion) TB() testing.TB { return a.tb }

// True 断言表达式 expr 为 true
//
// args 对应 fmt.Printf() 函数中的参数，其中 args[0] 对应第一个参数 format，依次类推，
// 具体可参数 Assert 方法的介绍。其它断言函数的 args 参数，功能与此相同。
func (a *Assertion) True(expr bool, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(expr, msg, []interface{}{"True 失败"})
}

func (a *Assertion) False(expr bool, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!expr, msg, []interface{}{"False 失败"})
}

func (a *Assertion) Nil(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(IsNil(expr), msg, []interface{}{"Nil 失败，实际值为 %#v", expr})
}

func (a *Assertion) NotNil(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!IsNil(expr), msg, []interface{}{"NotNil 失败，实际值为 %#v", expr})
}

func (a *Assertion) Equal(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(IsEqual(v1, v2), msg, []interface{}{"Equal 失败，实际值为\nv1=%#v\nv2=%#v", v1, v2})
}

func (a *Assertion) NotEqual(v1, v2 interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!IsEqual(v1, v2), msg, []interface{}{"NotEqual 失败，实际值为\nv1=%#v\nv2=%#v", v1, v2})
}

func (a *Assertion) Empty(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(IsEmpty(expr), msg, []interface{}{"Empty 失败，实际值为 %#v", expr})
}

func (a *Assertion) NotEmpty(expr interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!IsEmpty(expr), msg, []interface{}{"NotEmpty 失败，实际值为 %#v", expr})
}

// Error 断言有错误发生
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
//
// NotNil 的特化版本，限定了类型为 error。
func (a *Assertion) Error(expr error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!IsNil(expr), msg, []interface{}{"Error 失败，实际值为 Nil：%T", expr})
}

// ErrorString 断言有错误发生且错误信息中包含指定的字符串 str
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
func (a *Assertion) ErrorString(expr error, str string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if IsNil(expr) { // 空值，必定没有错误
		return a.Assert(false, msg, []interface{}{"ErrorString 失败，实际值为 Nil：%T", expr})
	}
	return a.Assert(strings.Contains(expr.Error(), str), msg, []interface{}{"ErrorString 失败，当前值为 %s", expr.Error()})
}

// ErrorIs 断言 expr 为 target 类型
//
// 相当于 a.True(errors.Is(expr, target))
func (a *Assertion) ErrorIs(expr, target error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(errors.Is(expr, target), msg, []interface{}{"ErrorIs 失败，expr 不是且不包含 target。"})
}

// NotError 断言没有错误
//
// Nil 的特化版本，限定了类型为 error。
func (a *Assertion) NotError(expr error, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(IsNil(expr), msg, []interface{}{"NotError 失败，实际值为：%s", expr})
}

func (a *Assertion) FileExists(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if _, err := os.Stat(path); err != nil && !os.IsExist(err) {
		return a.Assert(false, msg, []interface{}{"FileExists 失败，且附带以下错误：%v", err})
	}
	return a
}

func (a *Assertion) FileNotExists(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	_, err := os.Stat(path)
	if err == nil {
		return a.Assert(false, msg, []interface{}{"FileNotExists 失败"})
	}
	if os.IsExist(err) {
		return a.Assert(false, msg, []interface{}{"FileNotExists 失败，且返回以下错误信息：%v", err})
	}

	return a
}

func (a *Assertion) Panic(fn func(), msg ...interface{}) *Assertion {
	a.TB().Helper()

	has, _ := HasPanic(fn)
	return a.Assert(has, msg, []interface{}{"并未发生 panic"})
}

// PanicString 断言函数会发生 panic 且 panic 信息中包含指定的字符串内容
func (a *Assertion) PanicString(fn func(), str string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if has, m := HasPanic(fn); has {
		return a.Assert(strings.Contains(fmt.Sprint(m), str), msg, []interface{}{"panic 中并未包含 %s", str})
	}
	return a.Assert(false, msg, []interface{}{"并未发生 panic"})
}

// PanicType 断言函数会发生 panic 且抛出指定的类型
func (a *Assertion) PanicType(fn func(), typ interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if has, m := HasPanic(fn); has {
		t1, t2 := getType(true, m, typ)
		return a.Assert(t1 == t2, msg, []interface{}{"PanicType 失败，v1[%v] 的类型与 v2[%v] 的类型不相同", t1, t2})
	}
	return a.Assert(false, msg, []interface{}{"并未发生 panic"})
}

func (a *Assertion) NotPanic(fn func(), msg ...interface{}) *Assertion {
	a.TB().Helper()

	has, m := HasPanic(fn)
	return a.Assert(!has, msg, []interface{}{"发生了 panic，其信息为 %v", m})
}

// Contains 断言 container 包含 item 的或是包含 item 中的所有项
//
// 具体函数说明可参考 IsContains()
func (a *Assertion) Contains(container, item interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(IsContains(container, item), msg, []interface{}{"Contains 失败，%v 并未包含 %v", container, item})
}

// NotContains 断言 container 不包含 item 的或是不包含 item 中的所有项
func (a *Assertion) NotContains(container, item interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()
	return a.Assert(!IsContains(container, item), msg, []interface{}{"NotContains 失败，%v 包含 %v", container, item})
}

// Zero 断言是否为零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func (a *Assertion) Zero(v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()

	isZero := v == nil || reflect.ValueOf(v).IsZero()
	return a.Assert(isZero, msg, []interface{}{"%#v 为非零值", v})
}

// NotZero 断言是否为非零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func (a *Assertion) NotZero(v interface{}, msg ...interface{}) *Assertion {
	a.TB().Helper()

	isZero := v == nil || reflect.ValueOf(v).IsZero()
	return a.Assert(!isZero, msg, []interface{}{"%#v 为零值", v})
}

// Length 断言长度是否为指定的值
//
// v 可以是 map,string,slice,array
func (a *Assertion) Length(v interface{}, l int, msg ...interface{}) *Assertion {
	a.TB().Helper()

	rl, err := getLen(v)
	if err != "" {
		a.Assert(false, msg, []interface{}{err})
	}
	return a.Assert(rl == l, msg, []interface{}{"并非预期的长度，元素长度：%d, 期望的长度：%d", rl, l})
}

// NotLength 断言长度不是指定的值
//
// v 可以是 map,string,slice,array
func (a *Assertion) NotLength(v interface{}, l int, msg ...interface{}) *Assertion {
	a.TB().Helper()

	rl, err := getLen(v)
	if err != "" {
		a.Assert(false, msg, []interface{}{err})
	}
	return a.Assert(rl != l, msg, []interface{}{"长度均为 %d", rl})
}

// TypeEqual 断言两个值的类型是否相同
//
// ptr 如果为 true，则会在对象为指定时，查找其指向的对象。
func (a *Assertion) TypeEqual(ptr bool, v1, v2 interface{}, msg ...interface{}) *Assertion {
	if v1 == v2 {
		return a
	}

	a.TB().Helper()

	t1, t2 := getType(ptr, v1, v2)
	return a.Assert(t1 == t2, msg, []interface{}{"TypeEqual 失败，v1: %v，v2: %v", t1, t2})
}

func getType(ptr bool, v1, v2 interface{}) (t1, t2 reflect.Type) {
	t1 = reflect.TypeOf(v1)
	t2 = reflect.TypeOf(v2)

	if ptr {
		for t1.Kind() == reflect.Ptr {
			t1 = t1.Elem()
		}
		for t2.Kind() == reflect.Ptr {
			t2 = t2.Elem()
		}
	}

	return
}
