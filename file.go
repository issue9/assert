// SPDX-FileCopyrightText: 2014-2026 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"io/fs"
	"os"
)

func (a *Assertion) FileExists(path string, msg ...any) *Assertion {
	a.TB().Helper()
	_, err := os.Stat(path)
	return a.Assert(err == nil, NewFailure("FileExists", msg, map[string]any{"err": err}))
}

func (a *Assertion) FileNotExists(path string, msg ...any) *Assertion {
	a.TB().Helper()

	_, err := os.Stat(path)
	return a.Assert(errors.Is(err, fs.ErrNotExist), NewFailure("FileNotExists", msg, nil))
}

func (a *Assertion) FileExistsFS(fsys fs.FS, path string, msg ...any) *Assertion {
	a.TB().Helper()
	_, err := fs.Stat(fsys, path)
	return a.Assert(err == nil, NewFailure("FileExistsFS", msg, map[string]any{"err": err}))
}

func (a *Assertion) FileNotExistsFS(fsys fs.FS, path string, msg ...any) *Assertion {
	a.TB().Helper()

	_, err := fs.Stat(fsys, path)
	return a.Assert(errors.Is(err, fs.ErrNotExist), NewFailure("FileNotExistsFS", msg, nil))
}

// Dir 断言 path 是个目录
func (a *Assertion) Dir(path string, msg ...any) *Assertion {
	a.TB().Helper()

	s, err := os.Stat(path)
	if err != nil {
		return a.Assert(false, NewFailure("Dir", msg, map[string]any{"err": err}))
	}
	return a.Assert(s.IsDir(), NewFailure("Dir", msg, nil))
}

// NotDir 断言 path 不存在或是非目录
func (a *Assertion) NotDir(path string, msg ...any) *Assertion {
	a.TB().Helper()

	s, err := os.Stat(path)
	if err != nil {
		return a.Assert(true, NewFailure("NotDir", msg, map[string]any{"err": err}))
	}
	return a.Assert(!s.IsDir(), NewFailure("NotDir", msg, nil))
}

func (a *Assertion) DirFS(fsys fs.FS, path string, msg ...any) *Assertion {
	a.TB().Helper()

	s, err := fs.Stat(fsys, path)
	if err != nil {
		return a.Assert(false, NewFailure("DirFS", msg, map[string]any{"err": err}))
	}
	return a.Assert(s.IsDir(), NewFailure("DirFS", msg, nil))
}

func (a *Assertion) NotDirFS(fsys fs.FS, path string, msg ...any) *Assertion {
	a.TB().Helper()

	s, err := fs.Stat(fsys, path)
	if err != nil {
		return a.Assert(true, NewFailure("NotDirFS", msg, map[string]any{"err": err}))
	}
	return a.Assert(!s.IsDir(), NewFailure("NotDirFS", msg, nil))
}
