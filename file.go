package fsnotifyr

import (
	"io/fs"
	"os"
	"strings"
)

type File interface {
	fs.DirEntry
	FullPath() string
}

type file struct {
	folder Folder
	fs.DirEntry
}

func (f *file) FullPath() string {
	return strings.Join([]string{f.folder.FullPath(), f.DirEntry.Name()}, string(os.PathSeparator))
}

func NewFile(folder Folder, de fs.DirEntry) File {
	fyle := &file{folder, de}
	return fyle
}
