// Code generated by "esc -o static.go -prefix static static"; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return []os.FileInfo(fis[0:limit]), nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/index.html": {
		name:    "index.html",
		local:   "static/index.html",
		size:    5829,
		modtime: 1540421522,
		compressed: `
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
`,
	},

	"/": {
		name:  "/",
		local: `static`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"static": {
		_escData["/index.html"],
	},
}
