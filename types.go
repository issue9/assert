// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

type tester interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fail()
	FailNow()
	Failed() bool

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Log(args ...interface{})
	Logf(format string, args ...interface{})
}
