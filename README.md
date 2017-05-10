assert [![Build Status](https://travis-ci.org/issue9/assert.svg?branch=master)](https://travis-ci.org/issue9/assert)
======

assert包是对testing的一个简单扩展，提供的一系列的断言函数，
方便在测试函数中使用：
```go
func TestA(t testing.T) {
    v := true
    assert.True(v)

    a := assert.New(t)
    a.True(v)
}

// 也可以对testing.B使用
func Benchmark1(b *testing.B) {
    a := assert.New(b)
    v := false
    a.True(v)
    for(i:=0; i<b.N; i++) {
        // do something
    }
}
```

### 安装

```shell
go get github.com/issue9/assert
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/assert)
[![GoDoc](https://godoc.org/github.com/issue9/assert?status.svg)](https://godoc.org/github.com/issue9/assert)


### 版权

本项目采用[MIT](https://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
