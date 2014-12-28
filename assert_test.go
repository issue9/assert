// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"errors"
	"testing"
)

func TestGetCallerInfo(t *testing.T) {
	str := getCallerInfo()
	if len(str) == 0 {
		t.Error("getCallerInfo()无法正确返回信息")
	} else {
		t.Logf("getCallerInfo()返回的内容为：[%v]", str)
	}
}

func TestFormatMsg(t *testing.T) {
	msg1 := []interface{}{}
	msg2 := []interface{}{[]rune("msg:%v"), 2}
	msg3 := []interface{}{"msg:%v", 3}

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

func TestNotError(t *testing.T) {
	NotError(t, "123", "NotError(123) falid")

	var err1 error = nil
	NotError(t, err1, "var err1 error falid")

	err2 := &ErrorImpl{msg: "msg"}
	Error(t, err2, "Error(ErrorImpl) falid")
}

func TestFileExists(t *testing.T) {
	FileExists(t, "./assert.go", "FileExists() falid")
}

func TestFileNotExists(t *testing.T) {
	FileNotExists(t, "c:/win", "FileNotExists() falid")
}

func TestPanic(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	Panic(t, f1)
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

func TestStringEqual(t *testing.T) {
	StringEqual(t, "abc", "aBc", StyleCase)
}

func TestStringNotEqual(t *testing.T) {
	StringNotEqual(t, "abc", "aBc", StyleStrit)
}
