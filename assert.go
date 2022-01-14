// SPDX-License-Identifier: MIT

// Package assert 是对 testing 包的一些简单包装
//
//  func TestAssert(t *testing.T) {
//      var v interface{} = 5
//
//      a := assert.New(t, false)
//      a.True(v==5, "v的值[%v]不等于5", v).
//          Equal(5, v, "v的值[%v]不等于5", v).
//          Nil(v).
//          TB().Log("success")
//
//      // 以函数链的形式调用 Assertion 对象的方法
//      a.True(false).Equal(5,6)
//  }
//
//  // 也可以对 testing.B 使用
//  func Benchmark1(b *testing.B) {
//      a := assert.New(b)
//      a.True(false)
//      for(i:=0; i<b.N; i++) {
//          // do something
//      }
//  }
package assert
