package rwatch

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// a Folder contains zero or more files and zero or more subFolders
// it also contains a reference to the underlying filesystem mechanics
type Folder interface {
	fs.ReadDirFS
	SubFolders() []Folder
	ParentFolder() Folder // a pointer to the parent Folder, or nil for the King
	FullPath() string     // fully qualified path, taking ancestors into consideration
	//SpawnChild(string) (Folder, error)
	AddChild(Folder) error
	DestroyChild(Folder) error
	DestroySelf() error
	DestroyChildren() error
	Path() string // short path, relative to parent folder
}

// folder implements Folder
type folder struct {
	path       string
	parent     Folder
	filesystem fs.FS
	children   []Folder
}

func (f *folder) ReadDir(name string) ([]fs.DirEntry, error) {
	fullPath := strings.Join([]string{f.FullPath(), name}, string(os.PathSeparator))
	return fs.ReadDir(f.filesystem, fullPath)
}

func (f *folder) Open(name string) (fs.File, error) {
	return f.filesystem.Open(name)
}

// func (f *folder) Sub(name string) (fs.FS, error) {

// 	childFs, err := fs.Sub(f.filesystem, name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	child, err := NewFolder(childFs, name, f)

// 	return child, err
// }

func (f *folder) Path() string {
	return f.path
}

func (f *folder) ParentFolder() Folder {
	return f.parent
}

func (f *folder) Filesystem() fs.FS {
	return f.filesystem
}

// a subfolder is a Folder
func (f *folder) SpawnChild(subPath string) (Folder, error) {
	child, err := NewFolder(f.filesystem, subPath, f)
	if err != nil {
		err = f.AddChild(child)
	}
	return child, err
}

func (f *folder) SubFolders() []Folder {
	return f.children
}

func (f *folder) DestroyChildren() error {
	var err error
	for _, child := range f.children {
		err = child.DestroySelf()
		if err != nil {
			break
		}
	}
	if err != nil {
		f.children = nil
	}
	return err
}

// destroy all my children (which will act recursively) and then detroy myself
// by removing a reference to me from my parent
func (f *folder) DestroySelf() error {
	var err error
	err = f.DestroyChildren()
	if err == nil && f.parent != nil {
		err = f.parent.DestroyChild(f)
	}
	return err
}

func (f *folder) DestroyChild(thisChild Folder) error {
	err := fmt.Errorf("subfolder %q not found", thisChild.Path())
	for i, thatChild := range f.SubFolders() {
		if thisChild == thatChild {
			err = thisChild.DestroySelf()
			if err == nil {
				f.children = append(f.children[:i], f.children[i+1:]...)
			}
			break
		}
	}
	return err
}

// traverse up through ancestry to get full path
func (f *folder) FullPath() string {
	slugs := []string{f.path}
	visitor := f.parent
	for visitor != nil {
		slugs = append(slugs, visitor.Path())
		visitor = visitor.ParentFolder()
	}
	slices.Reverse(slugs)
	return filepath.Clean(strings.Join(slugs, string(os.PathListSeparator)))
}

func (f *folder) AddChild(child Folder) error {
	f.children = append(f.children, child)
	return nil
}

// This can be used to create the initial ancestor Folder, but also may be used as a helper function for Folder.Sub()
func NewFolder(filesystem fs.FS, path string, parent Folder) (Folder, error) {
	f := folder{
		path:       path,
		parent:     parent,
		filesystem: filesystem,
	}
	directoryContents, err := f.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s", path)
	}
	if parent != nil {
		parent.AddChild(&f)
	}

	//	immediately build recursive descendant tree
	for _, dirEntry := range directoryContents {
		switch dirEntry.Name() {
		case ".":
			if !dirEntry.IsDir() {
				return nil, fmt.Errorf("%s is not a directory", path)
			}
		default:
			if dirEntry.IsDir() {
				subFolder, err := NewFolder(f.filesystem, dirEntry.Name(), &f)
				if err != nil {
					return nil, err
				}
				f.children = append(f.children, subFolder)
			}
		}
	}

	return &f, nil
}
