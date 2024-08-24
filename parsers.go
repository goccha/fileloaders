package fileloaders

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type FileOption func(f *File)

func WithHash(hash *string) FileOption {
	return func(f *File) {
		f.hash = hash
	}
}

func WithVersion(version *string) FileOption {
	return func(f *File) {
		f.version = version
	}
}

type File struct {
	Type    string
	Bucket  string
	Path    string
	body    []byte
	hash    *string
	version *string
}

func (f *File) Hash() (string, bool) {
	if f.hash == nil {
		return "", false
	}
	return *f.hash, true
}

func (f *File) Version() (string, bool) {
	if f.version == nil {
		return "", false
	}
	return *f.version, true
}

func (f *File) GetBody() []byte {
	return f.body
}

func (f *File) Add(o ...FileOption) *File {
	for _, opt := range o {
		opt(f)
	}
	return f
}

func (f *File) WriteBody(p []byte) *File {
	f.body = p
	return f
}

func (f *File) Reader() io.ReadSeeker {
	return bytes.NewReader(f.body)
}

func (f *File) Write(p []byte) (n int, err error) {
	if f.body != nil {
		b := bytes.NewBuffer(f.body)
		b.Write(p)
		f.body = b.Bytes()
	} else {
		f.body = p
	}
	return len(p), nil
}
func (f *File) Unmarshal(obj any) error {
	return json.Unmarshal(f.body, obj)
}

func Parse(path string) *File {
	index := strings.Index(path, "://")
	if index < 0 {
		return nil
	}
	prefix := path[:index]
	path = path[index+3:]
	index = strings.Index(path, "/")
	var bucket string
	if index >= 0 {
		bucket = path[:index]
	}
	path = path[index+1:]
	return &File{
		Type:   prefix,
		Bucket: bucket,
		Path:   path,
	}
}
