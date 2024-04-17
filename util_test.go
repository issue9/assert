// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"testing"
	"time"
)

func TestIsZero(t *testing.T) {
	zero := func(v interface{}) {
		t.Helper()
		if !isZero(v) {
			t.Errorf("zero: %v", v)
		}
	}

	zero(nil)
	zero(struct{}{})
	zero(time.Time{})
	zero(&time.Time{})
}

func TestIsEqual(t *testing.T) {
	eq := func(v1, v2 interface{}) {
		t.Helper()
		if !isEqual(v1, v2) {
			t.Errorf("eq:[%v]!=[%v]", v1, v2)
		}
	}

	neq := func(v1, v2 interface{}) {
		t.Helper()
		if isEqual(v1, v2) {
			t.Errorf("eq:[%v]==[%v]", v1, v2)
		}
	}

	eq([]byte("abc"), "abc")
	eq("abc", []byte("abc"))

	eq([]byte("中文abc"), "中文abc")
	eq("中文abc", []byte("中文abc"))

	eq([]rune("中文abc"), "中文abc")
	eq("中文abc", []rune("中文abc"))

	eq(5, 5.0)
	eq(int8(5), 5)
	eq(5, int8(5))
	eq(float64(5), int8(5))
	eq([]int{1, 2, 3}, []int{1, 2, 3})
	eq([]int{1, 2, 3}, []int8{1, 2, 3})
	eq([]float32{1, 2.0, 3}, []int8{1, 2, 3})
	eq([]float32{1, 2.0, 3}, []float64{1, 2, 3})

	// 比较两个元素类型可相互转换的数组
	eq(
		[][]int{
			{1, 2},
			{3, 4},
		},
		[][]int8{
			{1, 2},
			{3, 4},
		},
	)

	// 比较两个元素类型可转换的 map
	eq(
		[]map[int]int{
			{1: 1, 2: 2},
			{3: 3, 4: 4},
		},
		[]map[int]int8{
			{1: 1, 2: 2},
			{3: 3, 4: 4},
		},
	)
	eq(map[string]int{"1": 1, "2": 2}, map[string]int8{"1": 1, "2": 2})

	// 比较两个元素类型可转换的 map
	eq(
		map[int]string{
			1: "1",
			2: "2",
		},
		map[int][]byte{
			1: []byte("1"),
			2: []byte("2"),
		},
	)

	// array 对比
	eq([2]int{1, 2}, [2]int{1, 2})
	eq([2]int{9, 3}, [2]int8{9, 3})
	eq([2]int8{1, 4}, [2]int{1, 4})
	eq([2]int{1, 5}, []int8{1, 5})

	neq(map[int]int{1: 1, 2: 2}, map[int8]int{1: 1, 2: 2})
	neq([]int{1, 2, 3}, []int{3, 2, 1})
	neq("5", 5)
	neq(true, "true")
	neq(true, 1)
	neq(true, "1")
	// 判断包含不同键名的两个 map
	neq(map[int]int{1: 1, 2: 2}, map[int]int{5: 5, 6: 6})

	// time
	loc := time.FixedZone("utf+8", 8*3600)
	now := time.Now()
	eq(time.Time{}, time.Time{})
	neq(now.In(loc), now.In(time.UTC)) // 时区不同
	n1 := time.Now()
	n2 := n1.Add(0)
	eq(n1, n2)

	// 指针
	v1 := 5
	v2 := 5
	p1 := &v1
	p2 := &v1
	eq(p1, p2) // 指针相等
	p2 = &v2
	eq(p1, p2) // 指向内容相等
}

func TestIsEmpty(t *testing.T) {
	if isEmpty([]string{""}) {
		t.Error("isEmpty([]string{\"\"})")
	}

	if !isEmpty([]string{}) {
		t.Error("isEmpty([]string{})")
	}

	if !isEmpty([]int{}) {
		t.Error("isEmpty([]int{})")
	}

	if !isEmpty(map[string]int{}) {
		t.Error("isEmpty(map[string]int{})")
	}

	if !isEmpty(0) {
		t.Error("isEmpty(0)")
	}

	if !isEmpty(int64(0)) {
		t.Error("isEmpty(int64(0))")
	}

	if !isEmpty(uint64(0)) {
		t.Error("isEmpty(uint64(0))")
	}

	if !isEmpty(0.0) {
		t.Error("isEmpty(0.0)")
	}

	if !isEmpty(float32(0)) {
		t.Error("isEmpty(0.0)")
	}

	if !isEmpty("") {
		t.Error("isEmpty(``)")
	}

	if !isEmpty([0]int{}) {
		t.Error("isEmpty([0]int{})")
	}

	if !isEmpty(time.Time{}) {
		t.Error("isEmpty(time.Time{})")
	}

	if !isEmpty(&time.Time{}) {
		t.Error("isEmpty(&time.Time{})")
	}

	if isEmpty("  ") {
		t.Error("isEmpty(\"  \")")
	}
}

func TestIsNil(t *testing.T) {
	if !isNil(nil) {
		t.Error("isNil(nil)")
	}

	var v1 []int
	if !isNil(v1) {
		t.Error("isNil(v1)")
	}

	var v2 map[string]string
	if !isNil(v2) {
		t.Error("isNil(v2)")
	}
}

func TestIsContains(t *testing.T) {
	fn := func(result bool, container, item interface{}) {
		t.Helper()
		if result != isContains(container, item) {
			t.Errorf("%v == (isContains(%v, %v))出错\n", result, container, item)
		}
	}

	fn(false, nil, nil)

	fn(true, "abc", "a")
	fn(true, "abc", "c")
	fn(true, "abc", "bc")
	fn(true, "abc", byte('a'))    // string vs byte
	fn(true, "abc", rune('a'))    // string vs rune
	fn(true, "abc", []byte("ab")) // string vs []byte
	fn(true, "abc", []rune("ab")) // string vs []rune

	fn(true, []byte("abc"), "a")
	fn(true, []byte("abc"), "c")
	fn(true, []byte("abc"), "bc")
	fn(true, []byte("abc"), byte('a'))
	fn(true, []byte("abc"), rune('a'))
	fn(true, []byte("abc"), []byte("ab"))
	fn(true, []byte("abc"), []rune("ab"))

	fn(true, []rune("abc"), "a")
	fn(true, []rune("abc"), "c")
	fn(true, []rune("abc"), "bc")
	fn(true, []rune("abc"), byte('a'))
	fn(true, []rune("abc"), rune('a'))
	fn(true, []rune("abc"), []byte("ab"))
	fn(true, []rune("abc"), []rune("ab"))

	fn(true, "中文a", "中")
	fn(true, "中文a", "a")
	fn(true, "中文a", '中')

	fn(true, []int{1, 2, 3}, 1)
	fn(true, []int{1, 2, 3}, int8(3))
	fn(true, []int{1, 2, 4}, []int{1, 2})
	fn(true, []interface{}{[]int{1, 2}, 5, 6}, []int8{1, 2})
	fn(true, []interface{}{[]int{1, 2}, 5, 6}, 5)

	fn(true, map[string]int{"1": 1, "2": 2}, map[string]int8{"1": 1})
	fn(true,
		map[string][]int{
			"1": {1, 2, 3},
			"2": {4, 5, 6},
		},
		map[string][]int8{
			"1": {1, 2, 3},
			"2": {4, 5, 6},
		},
	)

	fn(false, map[string]int{}, nil)
	fn(false, map[string]int{"1": 1, "2": 2}, map[string]int8{})
	fn(false, map[string]int{"1": 1, "2": 2}, map[string]int8{"1": 110}) // 同键名，不同值
	fn(false, map[string]int{"1": 1, "2": 2}, map[string]int8{"5": 5})
	fn(false, []int{1, 2, 3}, nil)
	fn(false, []int{1, 2, 3}, []int8{1, 3})
	fn(false, []int{1, 2, 3}, []int{1, 2, 3, 4})
	fn(false, []int{}, []int{1}) // 空数组
}
