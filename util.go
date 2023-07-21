// SPDX-License-Identifier: MIT

package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// 判断一个值是否为空(0, "", false, 空数组等)。
// []string{""}空数组里套一个空字符串，不会被判断为空。
func isEmpty(expr interface{}) bool {
	if isZero(expr) {
		return true
	}

	rv := reflect.ValueOf(expr)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan: // 长度为 0 的数组也是 empty
		return rv.Len() == 0
	default:
		return false
	}
}

func isZero(v interface{}) bool {
	if isNil(v) || reflect.ValueOf(v).IsZero() {
		return true
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.IsZero()
}

// isNil 判断一个值是否为 nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回 true
func isNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	return k >= reflect.Chan && k <= reflect.Slice && v.IsNil()
}

// 判断两个值是否相等
//
// 除了通过 reflect.DeepEqual() 判断值是否相等之外，一些类似
// 可转换的数值也能正确判断，比如以下值也将会被判断为相等：
//
//	int8(5)                     == int(5)
//	[]int{1,2}                  == []int8{1,2}
//	[]int{1,2}                  == [2]int8{1,2}
//	[]int{1,2}                  == []float32{1,2}
//	map[string]int{"1":"2":2}   == map[string]int8{"1":1,"2":2}
//
//	// map 的键值不同，即使可相互转换也判断不相等。
//	map[int]int{1:1,2:2}        != map[int8]int{1:1,2:2}
func isEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	vv1 := reflect.ValueOf(v1)
	vv2 := reflect.ValueOf(v2)

	if !vv1.IsValid() || !vv2.IsValid() {
		return vv1.IsValid() == vv2.IsValid()
	}

	if vv1 == vv2 {
		return true
	}

	vv1Type := vv1.Type()
	vv2Type := vv2.Type()

	if vv1Type.Comparable() && vv2Type.Comparable() && v1 == v2 {
		return true
	}

	// 过滤掉已经在 reflect.DeepEqual() 进行处理的类型
	switch vv1Type.Kind() {
	case reflect.Struct, reflect.Ptr, reflect.Func, reflect.Interface:
		return false
	case reflect.Slice, reflect.Array:
		// vv2.Kind() 与 vv1 的不相同
		if vv2.Kind() != reflect.Slice && vv2.Kind() != reflect.Array {
			// 虽然类型不同，但可以相互转换成 vv1 的，如：vv2 是 string，vv2 是 []byte，
			if vv2Type.ConvertibleTo(vv1Type) {
				return isEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
			}
			return false
		}

		// reflect.DeepEqual() 未考虑类型不同但是类型可转换的情况，比如：
		// []int{8,9} == []int8{8,9}，此处重新对 slice 和 array 做比较处理。
		if vv1.Len() != vv2.Len() {
			return false
		}

		for i := 0; i < vv1.Len(); i++ {
			if !isEqual(vv1.Index(i).Interface(), vv2.Index(i).Interface()) {
				return false
			}
		}
		return true // for 中所有的值比较都相等，返回 true
	case reflect.Map:
		if vv2.Kind() != reflect.Map {
			return false
		}

		if vv1.IsNil() != vv2.IsNil() {
			return false
		}
		if vv1.Len() != vv2.Len() {
			return false
		}
		if vv1.Pointer() == vv2.Pointer() {
			return true
		}

		// 两个 map 的键名类型不同
		if vv2Type.Key().Kind() != vv1Type.Key().Kind() {
			return false
		}

		for _, index := range vv1.MapKeys() {
			vv2Index := vv2.MapIndex(index)
			if !vv2Index.IsValid() {
				return false
			}

			if !isEqual(vv1.MapIndex(index).Interface(), vv2Index.Interface()) {
				return false
			}
		}
		return true // for 中所有的值比较都相等，返回 true
	case reflect.String:
		if vv2.Kind() == reflect.String {
			return vv1.String() == vv2.String()
		}
		if vv2Type.ConvertibleTo(vv1Type) { // 考虑 v1 是 string，v2 是 []byte 的情况
			return isEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
		}

		return false
	}

	if vv1Type.ConvertibleTo(vv2Type) {
		return vv2.Interface() == vv1.Convert(vv2Type).Interface()
	} else if vv2Type.ConvertibleTo(vv1Type) {
		return vv1.Interface() == vv2.Convert(vv1Type).Interface()
	}

	return false
}

// hasPanic 判断 fn 函数是否会发生 panic
// 若发生了 panic，将把 msg 一起返回。
func hasPanic(fn func()) (has bool, msg interface{}) {
	defer func() {
		if msg = recover(); msg != nil {
			has = true
		}
	}()
	fn()

	return
}

// isContains 判断 container 是否包含了 item 的内容。若是指针，会判断指针指向的内容
func isContains(container, item interface{}) bool {
	if container == nil { // nil不包含任何东西
		return false
	}

	cv := reflect.ValueOf(container)
	iv := reflect.ValueOf(item)

	for cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}

	for iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if isEqual(container, item) {
		return true
	}

	// 判断是字符串的情况
	switch c := cv.Interface().(type) {
	case string:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(c, i)
		case []byte:
			return strings.Contains(c, string(i))
		case []rune:
			return strings.Contains(c, string(i))
		case byte:
			return bytes.IndexByte([]byte(c), i) != -1
		case rune:
			return bytes.ContainsRune([]byte(c), i)
		}
	case []byte:
		switch i := iv.Interface().(type) {
		case string:
			return bytes.Contains(c, []byte(i))
		case []byte:
			return bytes.Contains(c, i)
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return bytes.IndexByte(c, i) != -1
		case rune:
			return bytes.ContainsRune(c, i)
		}
	case []rune:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(string(c), i)
		case []byte:
			return strings.Contains(string(c), string(i))
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return strings.IndexByte(string(c), i) != -1
		case rune:
			return strings.ContainsRune(string(c), i)
		}
	}

	if (cv.Kind() == reflect.Slice) || (cv.Kind() == reflect.Array) {
		if !cv.IsValid() || cv.Len() == 0 { // 空的，就不算包含另一个，即使另一个也是空值。
			return false
		}

		if !iv.IsValid() {
			return false
		}

		// item 是 container 的一个元素
		for i := 0; i < cv.Len(); i++ {
			if isEqual(cv.Index(i).Interface(), iv.Interface()) {
				return true
			}
		}

		// 开始判断 item 的元素是否与 container 中的元素相等。

		// 若 item 的长度为 0，表示不包含
		if (iv.Kind() != reflect.Slice) || (iv.Len() == 0) {
			return false
		}

		// item 的元素比 container 的元素多
		if iv.Len() > cv.Len() {
			return false
		}

		// 依次比较 item 的各个子元素是否都存在于 container，且下标都相同
		ivIndex := 0
		for i := 0; i < cv.Len(); i++ {
			if isEqual(cv.Index(i).Interface(), iv.Index(ivIndex).Interface()) {
				if (ivIndex == 0) && (i+iv.Len() > cv.Len()) {
					return false
				}
				ivIndex++
				if ivIndex == iv.Len() { // 已经遍历完 iv
					return true
				}
			} else if ivIndex > 0 {
				return false
			}
		}
		return false
	} // end cv.Kind == reflect.Slice and reflect.Array

	if cv.Kind() == reflect.Map {
		if cv.Len() == 0 {
			return false
		}

		if (iv.Kind() != reflect.Map) || (iv.Len() == 0) {
			return false
		}

		if iv.Len() > cv.Len() {
			return false
		}

		// 判断所有 item 的项都存在于 container 中
		for _, key := range iv.MapKeys() {
			cvItem := cv.MapIndex(key)
			if !cvItem.IsValid() { // container 中不包含该值。
				return false
			}
			if !isEqual(cvItem.Interface(), iv.MapIndex(key).Interface()) {
				return false
			}
		}
		// for 中的所有判断都成立，返回 true
		return true
	}

	return false
}

func getLen(v interface{}) (l int, msg string) {
	r := reflect.ValueOf(v)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	if v == nil {
		return 0, ""
	}

	switch r.Kind() {
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map, reflect.Chan:
		return r.Len(), ""
	}
	return 0, fmt.Sprintf("无法获取 %s 类型的长度信息", r.Kind())
}

func getType(ptr bool, v1, v2 interface{}) (t1, t2 reflect.Type) {
	t1 = reflect.TypeOf(v1)
	t2 = reflect.TypeOf(v2)

	if ptr {
		for t1.Kind() == reflect.Ptr {
			t1 = t1.Elem()
		}
		for t2.Kind() == reflect.Ptr {
			t2 = t2.Elem()
		}
	}

	return
}
