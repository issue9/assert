// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
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
	out := &bytes.Buffer{}
	a := newWithLogEnv(t, func(a ...any) { fmt.Fprint(out, a...) }, nil)

	if t != a.TB() {
		t.Error("a.T 与 t 不相等")
	}

	t.Run("True", func(t *testing.T) {
		out := &bytes.Buffer{}
		a := newWithLogEnv(t, func(a ...any) { fmt.Fprint(out, a...) }, nil)

		a.True(true).
			True(true, "a.True(5==5) failed")
		if out.Len() > 0 {
			t.Error("断言出错了")
		}

		a.True(false, "a.True(5==5) failed")
		if !strings.Contains(out.String(), "a.True(5==5) failed") {
			t.Error("断言出错了")
		}
	})

	t.Run("False", func(t *testing.T) {
		out := &bytes.Buffer{}
		a := newWithLogEnv(t, func(a ...any) { fmt.Fprint(out, a...) }, nil)

		a.False(false, "a.False(false) failed").
			False(false, "a.False(4==5) failed")
		if out.Len() > 0 {
			t.Error("断言出错了")
		}

		a.False(true, "a.False(true) failed")
		if !strings.Contains(out.String(), "a.False(true) failed") {
			t.Error("断言出错了")
		}
	})
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

	var v5 any
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

func TestAssertion_Zero_NotZero(t *testing.T) {
	a := New(t, false)

	var v any
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
}

func TestAssertion_Match(t *testing.T) {
	a := New(t, false)

	a.Match(regexp.MustCompile("^[1-9]*$"), "123")
	a.NotMatch(regexp.MustCompile("^[1-9]*$"), "x123")

	a.Match(regexp.MustCompile("^[1-9]*$"), []byte("123"))
	a.NotMatch(regexp.MustCompile("^[1-9]*$"), []byte("x123"))
}

func TestAssertion_When(t *testing.T) {
	a := New(t, false)

	a.When(true, func(a *Assertion) {
		a.True(true)
	})
}

func TestAssertion_Eventually(t *testing.T) {
	out := &bytes.Buffer{}
	a := newWithLogEnv(t, func(a ...any) { fmt.Fprint(out, a...) }, nil)

	cnt := 0
	a.Eventually(func() bool {
		cnt++
		if cnt > 10 {
			return true
		}
		return false
	}, 500*time.Microsecond)
	if out.Len() > 0 {
		t.Error("断言出错了")
	}

	a.Eventually(func() bool {
		return false
	}, 500*time.Microsecond, "Eventually always false")
	if !strings.Contains(out.String(), "Eventually always false") {
		t.Error("断言出错了")
	}
}

func TestAssertion_Never(t *testing.T) {
	out := &bytes.Buffer{}
	a := newWithLogEnv(t, func(a ...any) { fmt.Fprint(out, a...) }, nil)

	a.Never(func() bool { return false }, 500*time.Microsecond)
	if out.Len() > 0 {
		t.Error("断言出错了")
	}

	cnt := 0
	a.Never(func() bool {
		cnt++
		if cnt > 10 {
			return true
		}
		return false
	}, 500*time.Microsecond, "Never return true")
	if !strings.Contains(out.String(), "Never return true") {
		t.Error("断言出错了")
	}
}
