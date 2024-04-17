// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"database/sql"
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
