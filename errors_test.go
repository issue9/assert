// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestAssertion_Error(t *testing.T) {
	a := New(t, false)

	err := errors.New("test")
	a.Error(err, "a.Error(err) failed")
	a.ErrorString(err, "test", "ErrorString(err) failed")

	err2 := &errorImpl{msg: "msg"}
	a.Error(err2, "ErrorString(errorImpl) failed")
	a.ErrorString(err2, "msg", "ErrorString(errorImpl) failed")

	var err3 error
	a.NotError(err3, "var err1 error failed")

	err4 := errors.New("err4")
	err5 := fmt.Errorf("err5 with %w", err4)
	a.ErrorIs(err5, err4)
}

func TestAssertion_Panic(t *testing.T) {
	a := New(t, false)

	f1 := func() {
		panic("panic message")
	}

	a.Panic(f1)
	a.PanicString(f1, "message")
	a.PanicType(f1, "abc")
	a.PanicValue(f1, "panic message")

	f1 = func() {
		panic(errors.New("panic"))
	}
	a.PanicType(f1, errors.New("abc"))

	f1 = func() {
		panic(&errorImpl{msg: "panic"})
	}
	a.PanicType(f1, &errorImpl{msg: "abc"})

	f1 = func() {}
	a.NotPanic(f1)
}

func TestHasPanic(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	if has, _ := hasPanic(f1); !has {
		t.Error("f1未发生panic")
	}

	f2 := func() {
		f1()
	}

	if has, msg := hasPanic(f2); !has {
		t.Error("f2未发生panic")
	} else if msg != "panic" {
		t.Errorf("f2发生了panic，但返回信息不正确，应为[panic]，但其实返回了%v", msg)
	}

	f3 := func() {
		defer func() {
			if msg := recover(); msg != nil {
				t.Logf("TestHasPanic.f3 recover msg:[%v]", msg)
			}
		}()

		f1()
	}

	if has, msg := hasPanic(f3); has {
		t.Errorf("f3发生了panic，其信息为:[%v]", msg)
	}

	f4 := func() {
		//todo
	}

	if has, msg := hasPanic(f4); has {
		t.Errorf("f4发生panic，其信息为[%v]", msg)
	}
}
