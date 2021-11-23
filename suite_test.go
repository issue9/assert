// SPDX-License-Identifier: MIT

package assert

import "testing"

func TestTSuite(t *testing.T) {
	a := New(t)

	a.Run("r=1", func(a *Assertion) {
		a.True(true)
	}).Run("r=2", func(a *Assertion) {
		a.True(true)
		a.Run("r=2.1", func(a *Assertion) {
			a.Run("r=2.1.1", func(a *Assertion) {
				a.True(true)
			})
			a.True(true)
		})
	})
}
