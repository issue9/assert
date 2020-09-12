// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestAssertion(t *testing.T) {
	a := New(t)

	if t != a.TB() {
		t.Error("a.T与t不相等")
	}

	a.True(true)
	a.True(5 == 5, "a.True(5==5 failed")

	a.False(false, "a.False(false) failed")
	a.False(4 == 5, "a.False(4==5) failed")

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

	v9 := errors.New("test")
	a.Error(v9, "a.Error(v9) failed")

	a.NotError("abc", "a.NotError failed")

	a.FileExists("./assert.go", "a.FileExists(c:/windows) failed").
		FileNotExists("c:/win", "a.FileNotExists(c:/win) failed")

	err1 := errors.New("err1")
	err2 := fmt.Errorf("err2 with %w", err1)
	a.ErrorIs(err2, err1)
}
