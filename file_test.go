// SPDX-FileCopyrightText: 2014-2024 caixw
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

func TestAssertion_IsDir_IsNotDir(t *testing.T) {
	a := New(t, false)

	a.IsDir("./rest", "a.IsDir(./rest) failed").
		IsNotDir("./assert.go", "a.IsNotDir(./assert.go) failed")

	fsys := os.DirFS("./")
	a.IsDirFS(fsys, "rest", "a.IsDirFS(./rest) failed").
		IsNotDirFS(fsys, "./assert.go", "a.IsNotDirFS(./assert.go) failed")
}
