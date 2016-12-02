// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"errors"
	"testing"
)

func TestAssertion(t *testing.T) {
	a := New(t)

	if t != a.TB() {
		t.Error("a.T与t不相等")
	}

	a.True(true)
	a.True(5 == 5, "a.True(5==5 falid")

	a.False(false, "a.False(false) falid")
	a.False(4 == 5, "a.False(4==5) falid")

	v1 := 4
	v2 := 4
	v3 := 5
	v4 := "5"

	a.Equal(4, 4, "a.Equal(4,4) falid")
	a.Equal(v1, v2, "a.Equal(v1,v2) falid")

	a.NotEqual(4, 5, "a.NotEqual(4,5) falid").
		NotEqual(v1, v3, "a.NotEqual(v1,v3) falid").
		NotEqual(v3, v4, "a.NotEqual(v3,v4) falid")

	var v5 interface{}
	v6 := 0
	v7 := []int{}

	a.Empty(v5, "a.Empty falid").
		Empty(v6, "a.Empty(0) falid").
		Empty(v7, "a.Empty(v7) falid")

	a.NotEmpty(1, "a.NotEmpty(1) falid")

	a.Nil(v5)

	a.NotNil(v7, "a.Nil(v7) falid").
		NotNil(v6, "a.NotNil(v6) falid")

	v9 := errors.New("test")
	a.Error(v9, "a.Error(v9) falid")

	a.NotError("abc", "a.NotError falid")

	a.FileExists("./assert.go", "a.FileExists(c:/windows) falid").
		FileNotExists("c:/win", "a.FileNotExists(c:/win) falid")
}
