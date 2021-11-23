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

// 格式化错误提示信息
//
// msg1 输出的信息内容
// 所有参数将依次被传递给 fmt.Sprintf() 函数，
// 所以 msg1[0] 必须可以转换成 string(如:string, []byte, []rune, fmt.Stringer)
//
// msg2 参数格式与 msg1 完全相同，在 msg1 为空的情况下，会使用 msg2 的内容，
// 否则 msg2 不会启作用。
func formatMessage(msg1 []interface{}, msg2 []interface{}) string {
	if len(msg1) == 0 {
		msg1 = msg2
	}

	if len(msg1) == 0 {
		return "<未提供任何错误信息>"
	}

	if len(msg1) == 1 {
		return fmt.Sprint(msg1[0])
	}

	format := ""
	switch v := msg1[0].(type) {
	case []byte:
		format = string(v)
	case []rune:
		format = string(v)
	case string:
		format = v
	case fmt.Stringer:
		format = v.String()
	default:
		return fmt.Sprintln(msg1...)
	}

	return fmt.Sprintf(format, msg1[1:]...)
}

// Assert 断言 expr 条件成立
//
// expr 返回结果值为 bool 类型的表达式；
// msg1,msg2 输出的错误信息，之所以提供两组信息，是方便在用户没有提供的情况下，
// 可以使用系统内部提供的信息，优先使用 msg1 中的信息，若不存在，则使用 msg2 的内容。
//
// 直接使用 True 断言效果是一样的，之所以提供该函数，
// 主要供库调用，可以提供一个默认的错误信息。
func Assert(t testing.TB, expr bool, msg1 []interface{}, msg2 []interface{}) {
	t.Helper()
	if !expr {
		t.Error(formatMessage(msg1, msg2))
	}
}

// True 断言表达式 expr 为 true
//
// args 对应 fmt.Printf() 函数中的参数，其中 args[0] 对应第一个参数 format，依次类推，
// 具体可参数 formatMessage() 函数的介绍。其它断言函数的 args 参数，功能与此相同。
func True(t testing.TB, expr bool, args ...interface{}) {
	t.Helper()
	Assert(t, expr, args, []interface{}{"True 失败，实际值为 %#v", expr})
}

// False 断言表达式 expr 为 false
func False(t testing.TB, expr bool, args ...interface{}) {
	t.Helper()
	Assert(t, !expr, args, []interface{}{"False 失败，实际值为 %#v", expr})
}

// Nil 断言表达式 expr 为 nil
func Nil(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, IsNil(expr), args, []interface{}{"Nil 失败，实际值为 %#v", expr})
}

// NotNil 断言表达式 expr 为非 nil 值
func NotNil(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, !IsNil(expr), args, []interface{}{"NotNil 失败，实际值为 %#v", expr})
}

// Equal 断言 v1 与 v2 两个值相等
func Equal(t testing.TB, v1, v2 interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, IsEqual(v1, v2), args, []interface{}{"Equal 失败，实际值为\nv1=%#v\nv2=%#v", v1, v2})
}

// NotEqual 断言 v1 与 v2 两个值不相等
func NotEqual(t testing.TB, v1, v2 interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, !IsEqual(v1, v2), args, []interface{}{"NotEqual 失败，实际值为\nv1=%#v\nv2=%#v", v1, v2})
}

// Empty 断言 expr 的值为空(nil,"",0,false)，否则输出错误信息
func Empty(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, IsEmpty(expr), args, []interface{}{"Empty 失败，实际值为 %#v", expr})
}

// NotEmpty 断言 expr 的值为非空(除 nil,"",0,false之外)，否则输出错误信息
func NotEmpty(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()
	Assert(t, !IsEmpty(expr), args, []interface{}{"NotEmpty 失败，实际值为 %#v", expr})
}

// Error 断言有错误发生
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
func Error(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()

	if IsNil(expr) { // 空值，必定没有错误
		Assert(t, false, args, []interface{}{"Error 失败，实际值为 Nil：[%T]", expr})
		return
	}

	_, ok := expr.(error)
	Assert(t, ok, args, []interface{}{"Error 失败，实际类型为[%T]", expr})
}

// ErrorString 断言有错误发生且错误信息中包含指定的字符串 str
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败
func ErrorString(t testing.TB, expr interface{}, str string, args ...interface{}) {
	t.Helper()

	if IsNil(expr) { // 空值，必定没有错误
		Assert(t, false, args, []interface{}{"ErrorString 失败，实际值为 Nil：[%T]", expr})
		return
	}

	if err, ok := expr.(error); ok {
		index := strings.Index(err.Error(), str)
		Assert(t, index >= 0, args, []interface{}{"Error 失败，实际类型为[%T]", expr})
	}
}

// ErrorType 断言有错误发生且错误的类型与 typ 的类型相同
//
// 传递未初始化的 error 值(var err error = nil)，将断言失败。
//
// 仅对 expr 是否与 typ 为同一类型作简单判断，如果要检测是否是包含关系，可以使用 errors.Is 检测。
//
// ErrorType 与 ErrorIs 有本质的区别：ErrorIs 检测是否是包含关系，而 ErrorType 检测是否类型相同。比如：
//  err := os.WriteFile(...)
// 返回的 err 是一个 os.PathError 类型，用 ErrorType(err, &os.PathError{}) 断方正常；
// 而 ErrorIs(err, &os.PathError{}) 则会断言失败。
func ErrorType(t testing.TB, expr interface{}, typ error, args ...interface{}) {
	t.Helper()

	if IsNil(expr) { // 空值，必定没有错误
		Assert(t, false, args, []interface{}{"ErrorType 失败，实际值为 Nil：[%T]", expr})
		return
	}

	if _, ok := expr.(error); !ok {
		Assert(t, false, args, []interface{}{"ErrorType 失败，实际类型为[%T]，且无法转换成 error 接口", expr})
		return
	}

	t1 := reflect.TypeOf(expr)
	t2 := reflect.TypeOf(typ)
	Assert(t, t1 == t2, args, []interface{}{"ErrorType 失败，v1[%v]为一个错误类型，但与v2[%v]的类型不相同", t1, t2})
}

// NotError 断言没有错误发生
func NotError(t testing.TB, expr interface{}, args ...interface{}) {
	t.Helper()

	if IsNil(expr) { // 空值必定没有错误
		Assert(t, true, args, []interface{}{"NotError 失败，实际类型为[%T]", expr})
		return
	}
	err, ok := expr.(error)
	Assert(t, !ok, args, []interface{}{"NotError 失败，错误信息为[%v]", err})
}

// ErrorIs 断言 expr 为 target 类型
//
// 相当于 True(t, errors.Is(expr, target))
func ErrorIs(t testing.TB, expr interface{}, target error, args ...interface{}) {
	t.Helper()

	err, ok := expr.(error)
	Assert(t, ok, args, []interface{}{"ErrorIs 失败，expr 无法转换成 error。"})

	Assert(t, errors.Is(err, target), args, []interface{}{"ErrorIs 失败，expr 不是且不包含 target。"})
}

// FileExists 断言文件存在
func FileExists(t testing.TB, path string, args ...interface{}) {
	t.Helper()

	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		Assert(t, false, args, []interface{}{"FileExists 失败，且附带以下错误：%v", err})
	}
}

// FileNotExists 断言文件不存在
func FileNotExists(t testing.TB, path string, args ...interface{}) {
	t.Helper()

	_, err := os.Stat(path)
	if err == nil {
		Assert(t, false, args, []interface{}{"FileNotExists 失败"})
	}
	if os.IsExist(err) {
		Assert(t, false, args, []interface{}{"FileNotExists 失败，且返回以下错误信息：%v", err})
	}
}

// Panic 断言函数会发生 panic
func Panic(t testing.TB, fn func(), args ...interface{}) {
	t.Helper()

	has, _ := HasPanic(fn)
	Assert(t, has, args, []interface{}{"并未发生 panic"})
}

// PanicString 断言函数会发生 panic 且 panic 信息中包含指定的字符串内容
func PanicString(t testing.TB, fn func(), str string, args ...interface{}) {
	t.Helper()

	if has, msg := HasPanic(fn); has {
		index := strings.Index(fmt.Sprint(msg), str)
		Assert(t, index >= 0, args, []interface{}{"panic 中并未包含 %s", str})
		return
	}

	Assert(t, false, args, []interface{}{"并未发生 panic"})
}

// PanicType 断言函数会发生 panic 且抛出指定的类型
func PanicType(t testing.TB, fn func(), typ interface{}, args ...interface{}) {
	t.Helper()

	has, msg := HasPanic(fn)
	if !has {
		return
	}

	t1 := reflect.TypeOf(msg)
	t2 := reflect.TypeOf(typ)
	Assert(t, t1 == t2, args, []interface{}{"PanicType 失败，v1[%v]的类型与v2[%v]的类型不相同", t1, t2})

}

// NotPanic 断言函数不会发生 panic
func NotPanic(t testing.TB, fn func(), args ...interface{}) {
	t.Helper()

	has, msg := HasPanic(fn)
	Assert(t, !has, args, []interface{}{"发生了 panic，其信息为[%v]", msg})
}

// Contains 断言 container 包含 item 的或是包含 item 中的所有项
//
// 具体函数说明可参考 IsContains()
func Contains(t testing.TB, container, item interface{}, args ...interface{}) {
	t.Helper()

	Assert(t, IsContains(container, item), args,
		[]interface{}{"container:[%v]并未包含item[%v]", container, item})
}

// NotContains 断言 container 不包含 item 的或是不包含 item 中的所有项
func NotContains(t testing.TB, container, item interface{}, args ...interface{}) {
	t.Helper()

	Assert(t, !IsContains(container, item), args,
		[]interface{}{"container:[%v]包含item[%v]", container, item})
}

// Zero 断言是否为零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func Zero(t testing.TB, v interface{}, args ...interface{}) {
	t.Helper()

	isZero := v == nil || reflect.ValueOf(v).IsZero()
	Assert(t, isZero, args, []interface{}{"%#v 为非零值", v})
}

// NotZero 断言是否为非零值
//
// 最终调用的是 reflect.Value.IsZero 进行判断
func NotZero(t testing.TB, v interface{}, args ...interface{}) {
	t.Helper()

	isZero := v == nil || reflect.ValueOf(v).IsZero()
	Assert(t, !isZero, args, []interface{}{"%#v 为零值", v})
}

func Length(t testing.TB, v interface{}, l int, args ...interface{}) {
	t.Helper()

	rl := getLen(v)
	Assert(t, rl == l, args, []interface{}{"并非预期的长度，元素长度：%d, 期望的长度：%d", rl, l})
}

func NotLength(t testing.TB, v interface{}, l int, args ...interface{}) {
	t.Helper()

	rl := getLen(v)
	Assert(t, rl != l, args, []interface{}{"长度均为 %d", rl})
}

func getLen(v interface{}) int {
	r := reflect.ValueOf(v)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	switch r.Kind() {
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map:
		return r.Len()
	}
	return -1
}
