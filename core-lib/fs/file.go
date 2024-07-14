package fs

import (
	"io"
	"os"
	"path/filepath"
)

type File struct {
	path string
	info os.FileInfo
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) GetName() string {
	return f.info.Name()
}

func (f File) Extension() string {
	return filepath.Ext(f.path)
}

func (f File) Text() (string, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buff, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}

func NewFile(path string) (*File, error) {
	info, err := os.Lstat(path)
	if err == nil {
		file := File{path, info}
		return &file, nil
	}
	return nil, err
}
