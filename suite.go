// SPDX-License-Identifier: MIT

package assert

import "testing"

type (
	suite interface {
		run(name string, f func(a *Assertion))
	}

	tSuite struct {
		t *testing.T
	}

	bSuite struct {
		b *testing.B
	}
)

// Run 添加子测试
func (a *Assertion) Run(name string, f func(a *Assertion)) *Assertion {
	if a.suite == nil {
		switch obj := a.TB().(type) {
		case *testing.T:
			a.suite = &tSuite{t: obj}
		case *testing.B:
			a.suite = &bSuite{b: obj}
		default:
			panic("只有 *testing.T 和 *testing.B 支持 Run 功能")
		}
	}

	a.suite.run(name, f)
	return a
}

func (s *tSuite) run(name string, f func(a *Assertion)) {
	s.t.Run(name, func(t *testing.T) {
		f(New(t))
	})
}

func (s *bSuite) run(name string, f func(a *Assertion)) {
	s.b.Run(name, func(b *testing.B) {
		f(New(b))
	})
}
