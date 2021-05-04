package home

import (
	"fmt"
	"io/fs"
	"io/ioutil"
)

type IoUtil interface {
	ReadDir(dirname string) ([]fs.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

type Reader struct{}

func (r *Reader) ReadDir(dirname string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %s", err)
	}

	return files, nil
}

func (r *Reader) ReadFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}
