package fsnotifyr

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File interface {
	fs.DirEntry
	FullPath() string
	Parent() Folder
}

type file struct {
	folder Folder
	fs.DirEntry
}

func (f *file) Parent() Folder {
	return f.folder
}

func (f *file) FullPath() string {
	return filepath.Join(f.folder.FullPath(), f.DirEntry.Name())
}

func NewFile(folder Folder, de fs.DirEntry) File {
	return &file{folder, de}
}

func FileFromString(fullPath string, rootFolder Folder) File {

	parent := rootFolder
	basePath, _ := filepath.Split(fullPath)
	slugs := strings.Split(basePath, string(os.PathSeparator))

	fileInfo, err := fs.Stat(rootFolder.Filesystem(), fullPath)
	if err != nil {
		panic(err)
	}

slugs:
	for _, slug := range slugs {
		for _, child := range parent.Children() {
			if slug == child.Name() {
				parent = child
				continue slugs
			}
		}
	}
	return NewFile(parent, fs.FileInfoToDirEntry(fileInfo))

}
