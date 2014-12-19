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
```

### 安装

```shell
go get github.com/issue9/assert
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/issue9/assert)
[![GoDoc](https://godoc.org/github.com/issue9/assert?status.svg)](https://godoc.org/github.com/issue9/assert)


### 版权

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/issue9/assert/blob/master/LICENSE)
