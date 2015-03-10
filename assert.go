// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// 定位错误信息的触发函数。输出格式为：TestXxx(xxx_test.go:17)。
func getCallerInfo() string {
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		basename := path.Base(file)

		// 定位以_test.go结尾的文件。
		// 8 == len("_test.go")
		l := len(basename)
		if l < 8 || (basename[l-8:l] != "_test.go") {
			continue
		}

		// 定位函数名为Test开头的行。
		// 为什么要定位到TestXxx函数，是因为考虑以下情况：
		//  func isOK(val interface{}, t *testing.T) {
		//      // do somthing
		//      assert.True(t, val)  // (1
		//  }
		//
		//  func TestOK(t *testing.T) {
		//      isOK("123", t)       // (2
		//      isOK(123, t)         // (3
		//  }
		// 以上这段代码，定位到(2,(3的位置比总是定位到(1的位置更直观！
		funcName := runtime.FuncForPC(pc).Name()
		index := strings.LastIndex(funcName, ".Test")
		if -1 == index {
			continue
		}
		funcName = funcName[index+1:]
		if strings.IndexByte(funcName, '.') > -1 { // Go1.5之后的匿名函数
			continue
		}

		return funcName + "(" + basename + ":" + strconv.Itoa(line) + ")"
	}

	return "<无法获取调用者信息>"
}

// 格式化错误提示信息。
//
// msg1中的所有参数将依次被传递给fmt.Sprintf()函数，
// 所以msg1[0]必须可以转换成string(如:string, []byte, []rune, fmt.Stringer)
//
// msg2参数格式与msg1完全相同，在msg1为空的情况下，会使用msg2的内容，
// 否则msg2不会启作用。
func formatMessage(msg1 []interface{}, msg2 []interface{}) string {
	if len(msg1) == 0 {
		msg1 = msg2
	}

	if len(msg1) == 0 {
		return "<未提供任何错误信息>"
	}

	format := ""
	switch v := msg1[0].(type) {
	case []byte:
		format = string(v)
	case []rune:
		format = string(v)
	case string:
		format = v
	case fmt.Stringer:
		format = v.String()
	default:
		return "<无法正确转换错误提示信息>"
	}

	return fmt.Sprintf(format, msg1[1:]...)
}

// 当expr条件不成立时，输出错误信息。
//
// expr 返回结果值为bool类型的表达式；
// msg1,msg2输出的错误信息，之所以提供两组信息，是方便在用户没有提供的情况下，
// 可以使用系统内部提供的信息，优先使用msg1中的信息，若不存在，则使用msg2的内容。
func assert(t *testing.T, expr bool, msg1 []interface{}, msg2 []interface{}) {
	if !expr {
		t.Error(formatMessage(msg1, msg2) + "@" + getCallerInfo())
	}
}

// 断言表达式expr为true，否则输出错误信息。
//
// args对应fmt.Printf()函数中的参数，其中args[0]对应第一个参数format，依次类推，
// 具体可参数formatMessage()函数的介绍。其它断言函数的args参数，功能与此相同。
func True(t *testing.T, expr bool, args ...interface{}) {
	assert(t, expr, args, []interface{}{"True失败，实际值为[%T:%[1]v]", expr})
}

// 断言表达式expr为false，否则输出错误信息
func False(t *testing.T, expr bool, args ...interface{}) {
	assert(t, !expr, args, []interface{}{"False失败，实际值为[%T:%[1]v]", expr})
}

// 断言表达式expr为nil，否则输出错误信息
func Nil(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, IsNil(expr), args, []interface{}{"Nil失败，实际值为[%T:%[1]v]", expr})
}

// 断言表达式expr为非nil值，否则输出错误信息
func NotNil(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, !IsNil(expr), args, []interface{}{"NotNil失败，实际值为[%T:%[1]v]", expr})
}

// 断言v1与v2两个值相等，否则输出错误信息
func Equal(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	assert(t, IsEqual(v1, v2), args, []interface{}{"Equal失败，实际值为v1=[%T:%[1]v];v2=[%T:%[2]v]", v1, v2})
}

// 断言v1与v2两个值不相等，否则输出错误信息
func NotEqual(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	assert(t, !IsEqual(v1, v2), args, []interface{}{"NotEqual失败，实际值为v1=[%T:%[1]v];v2=[%T:%[2]v]", v1, v2})
}

// 断言expr的值为空(nil,"",0,false)，否则输出错误信息
func Empty(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, IsEmpty(expr), args, []interface{}{"Empty失败，实际值为[%T:%[1]v]", expr})
}

// 断言expr的值为非空(除nil,"",0,false之外)，否则输出错误信息
func NotEmpty(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, !IsEmpty(expr), args, []interface{}{"NotEmpty失败，实际值为[%T:%[1]v]", expr})
}

// 断言有错误发生，否则输出错误信息。
// 传递未初始化的error值(var err error = nil)，将断言失败
func Error(t *testing.T, expr interface{}, args ...interface{}) {
	if IsNil(expr) { // 空值，必定没有错误
		assert(t, false, args, []interface{}{"Error失败，实际类型为[%T]", expr})
	} else {
		_, ok := expr.(error)
		assert(t, ok, args, []interface{}{"Error失败，实际类型为[%T]", expr})
	}
}

// 断言没有错误发生，否则输出错误信息
func NotError(t *testing.T, expr interface{}, args ...interface{}) {
	if IsNil(expr) { // 空值必定没有错误
		assert(t, true, args, []interface{}{"NotError失败，实际类型为[%T]", expr})
	} else {
		err, ok := expr.(error)
		assert(t, !ok, args, []interface{}{"NotError失败，错误信息为[%v]", err})
	}
}

// 断言文件存在，否则输出错误信息
func FileExists(t *testing.T, path string, args ...interface{}) {
	_, err := os.Stat(path)

	if err != nil && !os.IsExist(err) {
		assert(t, false, args, []interface{}{"FileExists发生以下错误：%v", err.Error()})
	}
}

// 断言文件不存在，否则输出错误信息
func FileNotExists(t *testing.T, path string, args ...interface{}) {
	_, err := os.Stat(path)
	assert(t, os.IsNotExist(err), args, []interface{}{"FileExists发生以下错误：%v", err.Error()})
}

// 断言函数会发生panic，否则输出错误信息。
func Panic(t *testing.T, fn func(), args ...interface{}) {
	has, _ := HasPanic(fn)
	assert(t, has, args, []interface{}{"并未发生panic"})
}

// 断言函数会发生panic，否则输出错误信息。
func NotPanic(t *testing.T, fn func(), args ...interface{}) {
	has, msg := HasPanic(fn)
	assert(t, !has, args, []interface{}{"发生了panic，其信息为[%v]", msg})
}

// 断言container包含item的或是包含item中的所有项
// 具体函数说明可参考IsContains()
func Contains(t *testing.T, container, item interface{}, args ...interface{}) {
	assert(t, IsContains(container, item), args,
		[]interface{}{"container:[%v]并未包含item[%v]", container, item})
}

// 断言container不包含item的或是不包含item中的所有项
func NotContains(t *testing.T, container, item interface{}, args ...interface{}) {
	assert(t, !IsContains(container, item), args,
		[]interface{}{"container:[%v]包含item[%v]", container, item})
}
