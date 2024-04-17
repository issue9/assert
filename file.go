// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package assert

import (
	"errors"
	"io/fs"
	"os"
)

func (a *Assertion) FileExists(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if _, err := os.Stat(path); err != nil && !errors.Is(err, fs.ErrExist) {
		return a.Assert(false, NewFailure("FileExists", msg, map[string]interface{}{"err": err}))
	}
	return a
}

func (a *Assertion) FileNotExists(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	_, err := os.Stat(path)
	if err == nil {
		return a.Assert(false, NewFailure("FileNotExists", msg, nil))
	}
	if errors.Is(err, fs.ErrExist) {
		return a.Assert(false, NewFailure("FileNotExists", msg, map[string]interface{}{"err": err}))
	}

	return a
}

func (a *Assertion) FileExistsFS(fsys fs.FS, path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	if _, err := fs.Stat(fsys, path); err != nil && !errors.Is(err, fs.ErrExist) {
		return a.Assert(false, NewFailure("FileExistsFS", msg, map[string]interface{}{"err": err}))
	}

	return a
}

func (a *Assertion) FileNotExistsFS(fsys fs.FS, path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	_, err := fs.Stat(fsys, path)
	if err == nil {
		return a.Assert(false, NewFailure("FileNotExistsFS", msg, nil))
	}
	if errors.Is(err, fs.ErrExist) {
		return a.Assert(false, NewFailure("FileNotExistsFS", msg, map[string]interface{}{"err": err}))
	}

	return a
}

// IsDir 断言 path 是个目录
func (a *Assertion) IsDir(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	s, err := os.Stat(path)
	if err != nil {
		return a.Assert(false, NewFailure("IsDir", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(s.IsDir(), NewFailure("IsDir", msg, nil))
}

func (a *Assertion) IsDirFS(fsys fs.FS, path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	s, err := fs.Stat(fsys, path)
	if err != nil {
		return a.Assert(false, NewFailure("IsDirFS", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(s.IsDir(), NewFailure("IsDirFS", msg, nil))
}

// IsNotDir 断言 path 不存在或是非目录
func (a *Assertion) IsNotDir(path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	s, err := os.Stat(path)
	if err != nil {
		return a.Assert(false, NewFailure("IsNotDir", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(!s.IsDir(), NewFailure("IsNotDir", msg, nil))
}

func (a *Assertion) IsNotDirFS(fsys fs.FS, path string, msg ...interface{}) *Assertion {
	a.TB().Helper()

	s, err := os.Stat(path)
	if err != nil {
		return a.Assert(false, NewFailure("IsNotDirFS", msg, map[string]interface{}{"err": err}))
	}
	return a.Assert(!s.IsDir(), NewFailure("IsNotDirFS", msg, nil))
}
