// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package assert 测试用的断言包
//
//	func TestAssert(t *testing.T) {
//	    var v interface{} = 5
//
//	    a := assert.New(t, false)
//	    a.True(v==5, "v的值[%v]不等于5", v).
//	        Equal(5, v, "v的值[%v]不等于5", v).
//	        Nil(v).
//	        TB().Log("success")
//	}
//
//	// 也可以对 testing.B 使用
//	func Benchmark1(b *testing.B) {
//	    a := assert.New(b, false)
//	    a.True(false)
//	    for(i:=0; i<b.N; i++) {
//	        // do something
//	    }
//	}
package assert

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

var failureSprint FailureSprintFunc = DefaultFailureSprint

var failurePool = &sync.Pool{New: func() interface{} { return &Failure{} }}

// Failure 在断言出错时输出的错误信息
type Failure struct {
	Action string                 // 操作名称，比如 Equal，NotEqual 等方法名称。
	Values map[string]interface{} // 断言出错时返回的一些额外参数
	user   []interface{}          // 断言出错时用户反馈的额外信息
}

// FailureSprintFunc 将 [Failure] 转换成文本的函数
//
// NOTE: 可以使用此方法实现对错误信息的本地化。
type FailureSprintFunc = func(*Failure) string

// SetFailureSprintFunc 设置一个全局的转换方法
//
// [New] 方法在默认情况下继承由此方法设置的值。
func SetFailureSprintFunc(f FailureSprintFunc) { failureSprint = f }

// GetFailureSprintFunc 获取当前的 [FailureSprintFunc] 方法
func GetFailureSprintFunc() FailureSprintFunc { return failureSprint }

// DefaultFailureSprint 默认的 [FailureSprintFunc] 实现
func DefaultFailureSprint(f *Failure) string {
	s := strings.Builder{}
	s.WriteString(f.Action)
	s.WriteString(" 断言失败！")

	if len(f.Values) > 0 {
		keys := make([]string, 0, len(f.Values))
		for k := range f.Values {
			keys = append(keys, k)
		}
		sort.Strings(keys) // TODO(go1.21): slices.Sort

		s.WriteString("反馈以下参数：\n")
		for _, k := range keys {
			s.WriteString(k)
			s.WriteByte('=')
			s.WriteString(fmt.Sprint(f.Values[k]))
			s.WriteByte('\n')
		}
	}

	if u := f.User(); u != "" {
		s.WriteString("用户反馈信息：")
		s.WriteString(u)
	}

	return s.String()
}

// NewFailure 声明 [Failure] 对象
//
// user 表示用户提交的反馈，其第一个元素如果是 string，那么将调用 fmt.Sprintf(user[0], user[1:]...)
// 对数据进行格式化，否则采用 fmt.Sprint(user...) 格式化数据；
// kv 表示当前错误返回的数据；
func NewFailure(action string, user []interface{}, kv map[string]interface{}) *Failure {
	f := failurePool.Get().(*Failure)
	f.Action = action
	f.user = user
	f.Values = kv
	return f
}

// User 返回用户提交的返馈信息
func (f *Failure) User() string {
	// NOTE: 通过函数的方式返回字符串，而不是直接在 [NewFailure] 直接处理完，可以确保在未使用的情况下无需初始化。

	if len(f.user) == 0 {
		return ""
	}

	switch v := f.user[0].(type) {
	case string:
		return fmt.Sprintf(v, f.user[1:]...)
	default:
		return fmt.Sprint(f.user...)
	}
}
