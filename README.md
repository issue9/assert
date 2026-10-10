assert
======

[![Go](https://github.com/issue9/assert/workflows/Go/badge.svg)](https://github.com/issue9/assert/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/issue9/assert/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/assert)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/issue9/assert)](https://pkg.go.dev/github.com/issue9/assert/v5)
[![Go version](https://img.shields.io/github/go-mod/go-version/issue9/assert)](https://golang.org)

assert 包是对 testing 的一个简单扩展，提供的一系列的断言函数，方便在测试函数中使用。

```go
func TestA(t *testing.T) {
	v := true
	a := assert.New(t, false)
	a.True(v)
}

// 也可以对 testing.B 使用
func Benchmark1(b *testing.B) {
	a := assert.New(b, false)
	v := false
	a.True(v)
	for b.Loop() {
		// do something
	}
}
```

`assert/rest` 也提供了针对 HTTP 请求的测试功能

```go
func TestHTTP( t *testing.T) {
	a := assert.New(t, false)
	
	srv := rest.NewServer(a, h, nil)
	a.NotNil(srv)
	defer srv.Close()
	
	srv.NewRequest(http.MethodGet, "/body").
		Header("content-type", "application/json").
		Query("page", "5").
		EncodingBody(&bodyTest{ID: 5}, json.Marshal).
		Do().
		Status(http.StatusCreated).
		Header("content-type", "application/json;charset=utf-8").
		EncodingBody(&bodyTest{ID: 6}, json.Unmarshal)
}

// 也可以直接使用原始数据进行测试
func TestRaw(t *testing.T) {
	req :=`POST /users HTTP/1.1
	Host: example.com
	Content-type: application/json
	
	{"username": "admin", "password":"123"}
	
	`
	
	// 期望的返回数据
	resp :=`HTTP/1.1 201
	Location: https://example.com/users/1
	`

	a := assert.New(t, false)
	rest.RawHTTP(a, nil,req, resp)
}
```

版权
----

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
