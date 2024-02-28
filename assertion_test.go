// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

type errorImpl struct {
	msg string
}

func (err *errorImpl) Error() string {
	return err.msg
}

func TestAssertion_True_False(t *testing.T) {
	a := New(t, true)

	if t != a.TB() {
		t.Error("a.T与t不相等")
	}

	a.True(true)
	a.True(true, "a.True(5==5 failed")

	a.False(false, "a.False(false) failed")
	a.False(false, "a.False(4==5) failed")
}

func TestAssertion_Equal_NotEqual_Nil_NotNil(t *testing.T) {
	a := New(t, false)

	v1 := 4
	v2 := 4
	v3 := 5
	v4 := "5"

	a.Equal(4, 4, "a.Equal(4,4) failed")
	a.Equal(v1, v2, "a.Equal(v1,v2) failed")

	a.NotEqual(4, 5, "a.NotEqual(4,5) failed").
		NotEqual(v1, v3, "a.NotEqual(v1,v3) failed").
		NotEqual(v3, v4, "a.NotEqual(v3,v4) failed")

	var v5 interface{}
	v6 := 0
	v7 := []int{}

	a.Empty(v5, "a.Empty failed").
		Empty(v6, "a.Empty(0) failed").
		Empty(v7, "a.Empty(v7) failed")

	a.NotEmpty(1, "a.NotEmpty(1) failed")

	a.Nil(v5)

	a.NotNil(v7, "a.Nil(v7) failed").
		NotNil(v6, "a.NotNil(v6) failed")
}

func TestAssertion_Error(t *testing.T) {
	a := New(t, false)

	err := errors.New("test")
	a.Error(err, "a.Error(err) failed")
	a.ErrorString(err, "test", "ErrorString(err) failed")

	err2 := &errorImpl{msg: "msg"}
	a.Error(err2, "ErrorString(errorImpl) failed")
	a.ErrorString(err2, "msg", "ErrorString(errorImpl) failed")

	var err3 error
	a.NotError(err3, "var err1 error failed")

	err4 := errors.New("err4")
	err5 := fmt.Errorf("err5 with %w", err4)
	a.ErrorIs(err5, err4)
}

func TestAssertion_FileExists_FileNotExists(t *testing.T) {
	a := New(t, false)

	a.FileExists("./assert.go", "a.FileExists(c:/windows) failed").
		FileNotExists("c:/win", "a.FileNotExists(c:/win) failed")

	a.FileExistsFS(os.DirFS("./"), "assert.go", "a.FileExistsFS(c:/windows) failed").
		FileNotExistsFS(os.DirFS("c:/"), "win", "a.FileNotExistsFS(c:/win) failed")
}

func TestAssertion_Panic(t *testing.T) {
	a := New(t, false)

	f1 := func() {
		panic("panic message")
	}

	a.Panic(f1)
	a.PanicString(f1, "message")
	a.PanicType(f1, "abc")
	a.PanicValue(f1, "panic message")

	f1 = func() {
		panic(errors.New("panic"))
	}
	a.PanicType(f1, errors.New("abc"))

	f1 = func() {
		panic(&errorImpl{msg: "panic"})
	}
	a.PanicType(f1, &errorImpl{msg: "abc"})

	f1 = func() {}
	a.NotPanic(f1)
}

func TestAssertion_Zero_NotZero(t *testing.T) {
	a := New(t, false)

	var v interface{}
	a.Zero(0)
	a.Zero(nil)
	a.Zero(time.Time{})
	a.Zero(v)
	a.Zero([2]int{0, 0})
	a.Zero([0]int{})
	a.Zero(&time.Time{})
	a.Zero(sql.NullTime{})

	a.NotZero([]int{0, 0})
	a.NotZero([]int{})
}

func TestAssertion_Length_NotLength(t *testing.T) {
	a := New(t, false)

	a.Length(nil, 0)
	a.Length([]int{1, 2}, 2)
	a.Length([3]int{1, 2, 3}, 3)
	a.NotLength([3]int{1, 2, 3}, 2)
	a.Length(map[string]string{"1": "1", "2": "2"}, 2)
	a.NotLength(map[string]string{"1": "1", "2": "2"}, 3)
	slices := []rune{'a', 'b', 'c'}
	ps := &slices
	pps := &ps
	a.Length(pps, 3)
	a.NotLength(pps, 2)
	a.Length("string", 6)
	a.NotLength("string", 4)
}

func TestAssertion_Contains(t *testing.T) {
	a := New(t, false)

	a.Contains([]int{1, 2, 3}, []int8{1, 2}).
		NotContains([]int{1, 2, 3}, []int8{1, 3})
}

func TestAssertion_TypeEqual(t *testing.T) {
	a := New(t, true)

	a.TypeEqual(false, 1, 2)
	a.TypeEqual(false, 1, 1)
	a.TypeEqual(false, 1.0, 2.0)

	v1 := 5
	pv1 := &v1
	a.TypeEqual(false, 1, v1)
	a.TypeEqual(true, 1, &pv1)

	v2 := &errorImpl{}
	v3 := errorImpl{}
	a.TypeEqual(false, v2, v2)
	a.TypeEqual(true, v2, v3)
	a.TypeEqual(true, v2, &v3)
	a.TypeEqual(true, &v2, &v3)
}

func TestAssertion_Same(t *testing.T) {
	a := New(t, false)

	a.NotSame(5, 5).
		NotSame(struct{}{}, struct{}{}).
		NotSame(func() {}, func() {})

	i := 5
	a.NotSame(i, i)

	empty := struct{}{}
	empty2 := empty
	a.NotSame(empty, empty)
	a.NotSame(empty, empty2)
	a.Same(&empty, &empty)
	a.Same(&empty, &empty2)

	f := func() {}
	f2 := f
	a.Same(f, f)
	a.Same(f, f2)

	a.NotSame(5, 5)
	a.NotSame(f, 5)
}

func TestAssertion_Match(t *testing.T) {
	a := New(t, false)

	a.Match(regexp.MustCompile("^[1-9]*$"), "123")
	a.NotMatch(regexp.MustCompile("^[1-9]*$"), "x123")

	a.Match(regexp.MustCompile("^[1-9]*$"), []byte("123"))
	a.NotMatch(regexp.MustCompile("^[1-9]*$"), []byte("x123"))
}
