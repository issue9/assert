// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import "testing"

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

func TestAssertion_Greater_Less(t *testing.T) {
	a := New(t, false)

	a.Greater(uint16(5), 3).Less(uint8(5), 6).GreaterEqual(uint64(5), 5).LessEqual(uint(5), 5)
}

func TestAssertion_Positive_Negative(t *testing.T) {
	a := New(t, false)

	a.Positive(float32(5)).Negative(float64(-5))
}

func TestAssertion_Between(t *testing.T) {
	a := New(t, false)

	a.Between(int8(5), 1, 6).
		BetweenEqual(int16(5), 5, 6).
		BetweenEqual(int32(6), 5, 6).
		BetweenEqualMin(int64(5), 5, 6).
		BetweenEqualMax(uint32(5), 4, 5)
}
