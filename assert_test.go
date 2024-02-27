// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import "testing"

func TestDefaultFailureSprint(t *testing.T) {
	f := NewFailure("A", nil, nil)
	if f.Action != "A" || f.User != "" || len(f.Values) != 0 {
		t.Error("err1")
	}
	if s := DefaultFailureSprint(f); s != "A 断言失败！" {
		t.Error("err2")
	}

	// 带 user
	f = NewFailure("AB", []interface{}{1, 2}, nil)
	if f.Action != "AB" || f.User != "1 2" || len(f.Values) != 0 {
		t.Error("err3")
	}
	if s := DefaultFailureSprint(f); s != "AB 断言失败！用户反馈信息：1 2" {
		t.Error("err4", s)
	}

	// 带 values
	f = NewFailure("AB", nil, map[string]interface{}{"k1": "v1", "k2": 2})
	if f.Action != "AB" || f.User != "" || len(f.Values) != 2 {
		t.Error("err5")
	}
	if s := DefaultFailureSprint(f); s != "AB 断言失败！反馈以下参数：\nk1=v1\nk2=2\n" {
		t.Error("err6", s)
	}

	// 带 user,values
	f = NewFailure("AB", []interface{}{1, 2}, map[string]interface{}{"k1": "v1", "k2": 2})
	if f.Action != "AB" || f.User == "" || len(f.Values) != 2 {
		t.Error("err7")
	}
	if s := DefaultFailureSprint(f); s != "AB 断言失败！反馈以下参数：\nk1=v1\nk2=2\n用户反馈信息：1 2" {
		t.Error("err8", s)
	}
}
