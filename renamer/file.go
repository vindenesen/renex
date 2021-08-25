package renamer

import (
	"os"
	"path/filepath"
)

type File struct {
	directory     *Directory
	fileExtension string
	currentName   string
	originalName  string
	newName       string
}

func NewFile(directory *Directory, currentName string) *File {
	return &File{directory: directory, currentName: currentName, originalName: currentName}
}

func (f *File) CurrentName() string {
	return f.currentName
}

func (f *File) OriginalName() string {
	return f.originalName
}

func (f *File) SetNewName(newName string) string {
	f.newName = newName
	return f.NewName()
}

func (f *File) SetCurrentName(currentName string) string {
	f.currentName = currentName
	return f.CurrentName()
}

func (f *File) NewName() string {
	return f.newName
}

func (f *File) Directory() *Directory {
	return f.directory
}

func (f *File) FullPath() string {
	return filepath.Join(f.Directory().Path(), f.CurrentName())
}

func IsFile(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.Mode().IsRegular() {
		return true
	}
	return false
}
