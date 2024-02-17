package fsnotifyr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// a Folder contains zero or more files, and zero or more sub-Folders
// it also contains a reference to the underlying filesystem mechanics
type Folder interface {
	fs.ReadDirFS
	fs.DirEntry
	fmt.Stringer
	Children() []Folder
	Parent() Folder               // a pointer to the parent Folder, or nil for the root
	FullPath() string             // fully qualified path, prefixed using ancestors
	Spawn(string) (Folder, error) // creates a child Folder and adds it to the tree
	DestroyChild(Folder) error
	DestroyChildren() error
	Destroy() error
	FileTree(bool) FileTree
}

// folder implements Folder
type folder struct {
	path       string
	parent     Folder
	filesystem fs.FS
	children   []Folder
}

func (f *folder) Open(name string) (fs.File, error) {
	return f.filesystem.Open(name)
}

func (f *folder) String() string {
	// looks like the output of `tree -d` (directories only)
	// you can use FileTree().String() for output that includes regular files
	return f.FileTree(false).String()
}

func (f *folder) Stat() (fs.FileInfo, error) {
	return fs.Stat(f.filesystem, f.FullPath())
}

func (f *folder) Read(b []byte) (int, error) {
	return 0, fmt.Errorf("you can't read a directory as you would a file. %w", fs.ErrInvalid)
}

func (f *folder) Close() error {
	// what does it mean to close a directory?
	return nil
}

func (f *folder) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(f.filesystem, f.FullPath())
}

func (f *folder) Name() string {
	return f.path
}

func (f *folder) IsDir() bool {
	//	@note: This should always return true. a [Folder] after all is always a Dir.
	fd, err := f.filesystem.Open(f.FullPath())
	if err != nil {
		return false
	}
	info, err := fd.Stat()
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (f *folder) Info() (fs.FileInfo, error) {
	fd, err := f.filesystem.Open(f.FullPath())
	if err != nil {
		return nil, err
	}
	return fd.Stat()
}

func (f *folder) Type() fs.FileMode {
	fd, err := f.filesystem.Open(f.FullPath())
	if err != nil {
		return 0
	}
	info, err := fd.Stat()
	if err != nil {
		return 0
	}
	return info.Mode().Type()
}

func (f *folder) Path() string {
	return f.path
}

func (f *folder) ParentFolder() Folder {
	return f.parent
}

func (f *folder) Filesystem() fs.FS {
	return f.filesystem
}

// Spawn creates a new Folder as a child of this Folder
// It also recursively
func (f *folder) Spawn(subPath string) (Folder, error) {
	child, err := NewFolder(f.filesystem, subPath, f)
	return child, err
}

func (f *folder) Parent() Folder {
	return f.parent
}

func (f *folder) Children() []Folder {
	return f.children
}

func (f *folder) Destroy() error {
	if f.Parent() != nil {
		return f.Parent().DestroyChild(f)
	}
	//	@todo: what does it mean to destroy the root folder?
	return nil
}

func (f *folder) DestroyChildren() error {
	var err error
	for _, child := range f.children {
		err = child.DestroyChildren()
		if err != nil {
			break
		}
	}
	if err != nil {
		f.children = nil
	}
	return err
}

func (f *folder) DestroyChild(thisChild Folder) error {
	err := fmt.Errorf("subfolder %q not found", thisChild.Name())
	for i, thatChild := range f.Children() {
		if thisChild == thatChild {
			err = thisChild.DestroyChildren()
			if err == nil {
				f.children = append(f.children[:i], f.children[i+1:]...)
			}
			break
		}
	}
	return err
}

// Traverse up through ancestry to get full path
// This is not the necisarily the _absolute_ path
// but it allows children to know where they are in the hierarchy
// relative to some root that makes sense to the [fs.FS] instance passed in
// such as a filesystem's "current working directory"
func (f *folder) FullPath() string {
	slugs := []string{}
	visitor := f
	for visitor != nil {
		slugs = append(slugs, visitor.Name())
		//		}
		p := visitor.Parent()
		if p != nil {
			visitor = p.(*folder)
		} else {
			visitor = nil
		}
	}
	slices.Reverse(slugs)
	return filepath.Clean(strings.Join(slugs, string(os.PathSeparator)))
}

func (f *folder) addChild(child Folder) error {
	f.children = append(f.children, child)
	return nil
}

// NewFolder is the constructor for [Folder]. It also builds out its descendent tree
func NewFolder(filesystem fs.FS, path string, parent Folder) (Folder, error) {
	f := folder{
		path:       path,
		parent:     parent,
		filesystem: filesystem,
	}
	directoryContents, err := f.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %q", path)
	}
	if parent != nil {
		parent.(*folder).addChild(&f)
	}

	//	build recursive descendant tree
	for _, dirEntry := range directoryContents {
		switch dirEntry.Name() {
		case ".":
			if !dirEntry.IsDir() {
				return nil, fmt.Errorf("%s is not a directory", path)
			}
		default:
			if dirEntry.IsDir() {
				_, err := NewFolder(f.filesystem, dirEntry.Name(), &f)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &f, nil
}
