// SPDX-License-Identifier: MIT

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
	if str != "TestGetCallerInfo(assert_test.go:22)" {
		t.Errorf("getCallerInfo 返回的信息不正确，其返回值为：%v", str)
	}

	// 嵌套调用，第二个参数为当前的行号
	testGetCallerInfo(t, "29")
	testGetCallerInfo(t, "30")

	// 闭合函数，line 为调用所在的行号。
	f := func(line string) {
		str := getCallerInfo()
		if str != "TestGetCallerInfo(assert_test.go:"+line+")" {
			t.Errorf("getCallerInfo 返回的信息不正确，其返回值为：%v", str)
		}
	}

	go func() {
		f("41")
	}()
	go func() {
		testGetCallerInfo(t, "44")
	}()

	// bug: 无法处理的情况，go 会新开协程，无法获取当前的行号
	//go f("50")
	//go testGetCallerInfo(t, "51")

	f("51") // 参数为当前等号
	f("52")

	ff := func(line string) {
		f(line)
	}
	go func() {
		ff("58")
	}()

	// NOTE: 无法处理的情况
	/*go func() {
		go ff("63")
	}()*/

	time.Sleep(500 * time.Microsecond)
}

// 参数 line，为调用此函数所在的行号。
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
	True(t, 1 == 1, "True(1==1) failed")
}

func TestFalse(t *testing.T) {
	False(t, false, "False failed")
	False(t, 1 == 2, "False(1==2) failed")
}

func TestNil(t *testing.T) {
	Nil(t, nil, "Nil failed")

	var v interface{}
	Nil(t, v, "Nil(v) failed")
}

func TestNotNil(t *testing.T) {
	NotNil(t, 5, "NotNil failed")

	var v interface{} = 5
	NotNil(t, v, "NotNil failed")
}

func TestEqual(t *testing.T) {
	Equal(t, 5, 5, "Equal(5,5) failed")

	var v1, v2 interface{}
	v1 = 5
	v2 = 5

	Equal(t, 5, v1)
	Equal(t, v1, v2, "Equal(v1,v2) failed")
	Equal(t, int8(126), 126)
	Equal(t, int64(126), int8(126))
	Equal(t, uint(7), int(7))
}

func TestNotEqual(t *testing.T) {
	NotEqual(t, 5, 6, "NotEqual(5,6) failed")

	var v1, v2 interface{} = 5, 6

	NotEqual(t, 5, v2, "NotEqual(5,v2) failed")
	NotEqual(t, v1, v2, "NotEqual(v1,v2) failed")
	NotEqual(t, 128, int8(127))
}

func TestEmpty(t *testing.T) {
	Empty(t, 0, "Empty(0) failed")
	Empty(t, "", "Empty(``) failed")
	Empty(t, false, "Empty(false) failed")
	Empty(t, []string{}, "Empty(slice{}) failed")
	Empty(t, []int{}, "Empty(slice{}) failed")
}

func TestNotEmpty(t *testing.T) {
	NotEmpty(t, 1, "NotEmpty(1) failed")
	NotEmpty(t, true, "NotEmpty(true) failed")
	NotEmpty(t, []string{"ab"}, "NotEmpty(slice(abc)) failed")
}

type ErrorImpl struct {
	msg string
}

func (err *ErrorImpl) Error() string {
	return err.msg
}

func TestError(t *testing.T) {
	err1 := errors.New("test")
	Error(t, err1, "Error(err) failed")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "Error(ErrorImpl) failed")
}

func TestErrorString(t *testing.T) {
	err1 := errors.New("test")
	ErrorString(t, err1, "test", "Error(err1) failed")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "msg", "Error(ErrorImpl) failed")
}

func TestErrorType(t *testing.T) {
	ErrorType(t, errors.New("abc"), errors.New("def"), "ErrorType:errors.New(abc) != errors.New(def)")

	ErrorType(t, &ErrorImpl{msg: "abc"}, &ErrorImpl{}, "ErrorType:&ErrorImpl{} != &ErrorImpl{}")
}

func TestNotError(t *testing.T) {
	NotError(t, "123", "NotError(123) failed")

	var err1 error
	NotError(t, err1, "var err1 error failed")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "Error(ErrorImpl) failed")
}

func TestFileExists(t *testing.T) {
	FileExists(t, "./assert.go", "FileExists() failed")
	FileExists(t, "./", "FileExists() failed")
}

func TestFileNotExists(t *testing.T) {
	FileNotExists(t, "c:/win", "FileNotExists() failed")
	FileNotExists(t, "./abcefg/", "FileNotExists() failed")
}

func TestPanic(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	Panic(t, f1)
}

func TestPanicString(t *testing.T) {
	f1 := func() {
		panic("panic message")
	}

	PanicString(t, f1, "panic message")
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
	f1 := func() {}

	NotPanic(t, f1)
}

func TestContains(t *testing.T) {
	Contains(t, []int{1, 2, 3}, []int8{1, 2})
}

func TestNotContains(t *testing.T) {
	NotContains(t, []int{1, 2, 3}, []int8{1, 3})
}

func TestZero(t *testing.T) {
	Zero(t, 0)
	Zero(t, time.Time{})
	NotZero(t, &time.Time{})
}
