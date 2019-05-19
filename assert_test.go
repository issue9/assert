// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"errors"
	"testing"
	"time"
)

func BenchmarkGetCallerInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := getCallerInfo()

		if str != "BenchmarkGetCallerInfo(assert_test.go:15)" {
			b.Errorf("getCallerInfo 返回的信息不正确，其返回值为：%v", str)
		}
	}
}

func TestGetCallerInfo(t *testing.T) {
	str := getCallerInfo()
	// NOTE:注意这里涉及到调用函数的行号信息
	if str != "TestGetCallerInfo(assert_test.go:24)" {
		t.Errorf("getCallerInfo返回的信息不正确，其返回值为：%v", str)
	}

	// 嵌套调用，第二个参数为当前的行号
	testGetCallerInfo(t, "31")
	testGetCallerInfo(t, "32")

	// 闭合函数，line 为调用所在的行号。
	f := func(line string) {
		str := getCallerInfo()
		if str != "TestGetCallerInfo(assert_test.go:"+line+")" {
			t.Errorf("getCallerInfo返回的信息不正确，其返回值为：%v", str)
		}
	}

	go func() {
		f("43")
	}()
	go func() {
		testGetCallerInfo(t, "46")
	}()

	// NOTE: 无法处理的情况
	//go f("49")
	//go testGetCallerInfo(t, "50")

	f("53") // 参数为当前等号
	f("54")

	ff := func(line string) {
		f(line)
	}
	go func() {
		ff("60")
	}()

	// NOTE: 无法处理的情况
	/*go func() {
		go ff("63")
	}()*/

	time.Sleep(500 * time.Microsecond)
}

// 参数line，为调用此函数所在的行号。
func testGetCallerInfo(t *testing.T, line string) {
	str := getCallerInfo()
	if str != "TestGetCallerInfo(assert_test.go:"+line+")" {
		t.Errorf("getCallerInfo返回的信息不正确，其返回值为：%v", str)
	}
}

func TestFormatMsg(t *testing.T) {
	msg1 := []interface{}{}
	msg2 := []interface{}{[]rune("msg:%v"), 2}
	msg3 := []interface{}{"msg:%v", 3}
	msg4 := []interface{}{123, 456}
	msg5 := []interface{}{123}
	msg6 := []interface{}{errors.New("123")}

	str := formatMessage(msg1, msg2)
	if str != "msg:2" {
		t.Errorf("formatMessage(msg1,msg2)返回信息错误:[%v]", str)
	}

	str = formatMessage(nil, msg2)
	if str != "msg:2" {
		t.Errorf("formatMessage(msg1,msg2)返回信息错误:[%v]", str)
	}

	str = formatMessage(msg2, msg3)
	if str != "msg:2" {
		t.Errorf("formatMessage(msg2,msg3)返回信息错误:[%v]", str)
	}

	str = formatMessage(nil, nil)
	if str != "<未提供任何错误信息>" {
		t.Errorf("formatMessage(nil,nil)返回信息错误:[%v]", str)
	}

	str = formatMessage(nil, msg4)
	if str != "123 456\n" {
		t.Errorf("formatMessage(nil,nil)返回信息错误:[%v]", str)
	}

	str = formatMessage(nil, msg5)
	if str != "123" {
		t.Errorf("formatMessage(nil,msg5)返回信息错误:[%v]", str)
	}

	str = formatMessage(nil, msg6)
	if str != "123" {
		t.Errorf("formatMessage(nil,msg5)返回信息错误:[%v]", str)
	}
}

func TestTrue(t *testing.T) {
	True(t, true)
	True(t, 1 == 1, "True(1==1) falid")
}

func TestFalse(t *testing.T) {
	False(t, false, "False falid")
	False(t, 1 == 2, "False(1==2) falid")
}

func TestNil(t *testing.T) {
	Nil(t, nil, "Nil falid")

	var v interface{}
	Nil(t, v, "Nil(v) falid")
}

func TestNotNil(t *testing.T) {
	NotNil(t, 5, "NotNil falid")

	var v interface{} = 5
	NotNil(t, v, "NotNil falid")
}

func TestEqual(t *testing.T) {
	Equal(t, 5, 5, "Equal(5,5) falid")

	var v1, v2 interface{}
	v1 = 5
	v2 = 5

	Equal(t, 5, v1)
	Equal(t, v1, v2, "Equal(v1,v2) falid")
	Equal(t, int8(126), 126)
	Equal(t, int64(126), int8(126))
	Equal(t, uint(7), int(7))
}

func TestNotEqual(t *testing.T) {
	NotEqual(t, 5, 6, "NotEqual(5,6) falid")

	var v1, v2 interface{} = 5, 6

	NotEqual(t, 5, v2, "NotEqual(5,v2) falid")
	NotEqual(t, v1, v2, "NotEqual(v1,v2) falid")
	NotEqual(t, 128, int8(127))
}

func TestEmpty(t *testing.T) {
	Empty(t, 0, "Empty(0) falid")
	Empty(t, "", "Empty(``) falid")
	Empty(t, false, "Empty(false) falid")
	Empty(t, []string{}, "Empty(slice{}) falid")
	Empty(t, []int{}, "Empty(slice{}) falid")
}

func TestNotEmpty(t *testing.T) {
	NotEmpty(t, 1, "NotEmpty(1) falid")
	NotEmpty(t, true, "NotEmpty(true) falid")
	NotEmpty(t, []string{"ab"}, "NotEmpty(slice(abc)) falid")
}

type ErrorImpl struct {
	msg string
}

func (err *ErrorImpl) Error() string {
	return err.msg
}

func TestError(t *testing.T) {
	err1 := errors.New("test")
	Error(t, err1, "Error(err) falid")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "Error(ErrorImpl) falid")
}

func TestErrorString(t *testing.T) {
	err1 := errors.New("test")
	ErrorString(t, err1, "test", "Error(err1) falid")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "msg", "Error(ErrorImpl) falid")
}

func TestErrorType(t *testing.T) {
	ErrorType(t, errors.New("abc"), errors.New("def"), "ErrorType:errors.New(abc) != errors.New(def)")

	ErrorType(t, &ErrorImpl{msg: "abc"}, &ErrorImpl{}, "ErrorType:&ErrorImpl{} != &ErrorImpl{}")
}

func TestNotError(t *testing.T) {
	NotError(t, "123", "NotError(123) falid")

	var err1 error
	NotError(t, err1, "var err1 error falid")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "Error(ErrorImpl) falid")
}

func TestFileExists(t *testing.T) {
	FileExists(t, "./assert.go", "FileExists() falid")
	FileExists(t, "./", "FileExists() falid")
}

func TestFileNotExists(t *testing.T) {
	FileNotExists(t, "c:/win", "FileNotExists() falid")
	FileNotExists(t, "./abcefg/", "FileNotExists() falid")
}

func TestPanic(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	Panic(t, f1)
}

func TestPanicString(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	PanicString(t, f1, "pani")
}

func TestPanicType(t *testing.T) {
	f1 := func() {
		panic("panic")
	}
	PanicType(t, f1, "abc")

	f1 = func() {
		panic(errors.New("panic"))
	}
	PanicType(t, f1, errors.New("abc"))

	f1 = func() {
		panic(&ErrorImpl{msg: "panic"})
	}
	PanicType(t, f1, &ErrorImpl{msg: "abc"})
}

func TestNotPanic(t *testing.T) {
	f1 := func() {
	}

	NotPanic(t, f1)
}

func TestContains(t *testing.T) {
	Contains(t, []int{1, 2, 3}, []int8{1, 2})
}

func TestNotContains(t *testing.T) {
	NotContains(t, []int{1, 2, 3}, []int8{1, 3})
}
