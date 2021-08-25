package renamer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Directory struct {
	path  string
	files *[]*File
}

func NewDirectory(directory string) (*Directory, error) {
	if IsDirectory(directory) {
		newDir := &Directory{path: directory}
		err := newDir.findMyFiles()
		if err != nil {
			return nil, err
		}
		return newDir, nil
	}
	return nil, errors.New("directory doesn't exist or is inaccessible")
}

func (d *Directory) Path() string {
	return d.path
}

func (d *Directory) Files() *[]*File {
	return d.files
}

func (d *Directory) findMyFiles() error {
	content, err := os.ReadDir(d.Path())
	if err != nil {
		return errors.New(fmt.Sprintf("failed during reading of directory %s: %s", d.path, err.Error()))
	}

	files := make([]*File, 0)
	for _, f := range content {
		files = append(files, NewFile(d, f.Name()))
	}
	d.files = &files
	return nil
}

func IsDirectory(dir string) bool {
	dir = filepath.Clean(dir)
	fi, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	if fi.Mode().IsDir() {
		return true
	}
	return false
}
