// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"os"
	"testing"
)

func TestAssertion_FileExists_FileNotExists(t *testing.T) {
	a := New(t, false)

	a.FileExists("./assert.go", "a.FileExists(./assert.go) failed").
		FileNotExists("c:/win", "a.FileNotExists(c:/win) failed")

	fsys := os.DirFS("./")
	a.FileExistsFS(fsys, "assert.go", "a.FileExistsFS(./assert) failed").
		FileNotExistsFS(fsys, "win", "a.FileNotExistsFS(c:/win) failed")
}

func TestAssertion_Dir_NotDir(t *testing.T) {
	a := New(t, false)

	a.Dir("./rest", "a.Dir(./rest) failed").
		NotDir("./assert.go", "a.NotDir(./assert.go) failed")

	fsys := os.DirFS("./")
	a.DirFS(fsys, "rest", "a.DirFS(./rest) failed").
		NotDirFS(fsys, "assert.go", "a.NotDirFS(./assert.go) failed")
}
