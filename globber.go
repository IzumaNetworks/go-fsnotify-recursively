package fsnotifyr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

/**
 *	Globber seperates out the non-magical part of a glob string, compiles a glob, and retains all necessary information
 */

type Globber interface {
	fmt.Stringer
	Root() string
	Match(str string) bool
}

type globber struct {
	fsRoot  string
	pattern string
}

func NewGlobber(fullString string) (Globber, error) {
	fsRoot, globRoot, err := ComponentizeGlobString(fullString)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	globber := &globber{fsRoot, globRoot}
	return globber, nil
}

func (g *globber) Match(str string) bool {
	matches, err := doublestar.Match(g.pattern, str)
	if err != nil {
		panic(err)
	}
	return matches
}

func (g *globber) Root() string {
	return g.fsRoot
}

func (g *globber) String() string {
	output := map[string]string{
		"fsRoot":   g.fsRoot,
		"globRoot": g.pattern,
	}
	j, _ := json.Marshal(output)
	return string(j)
}

func ComponentizeGlobString(globExpression string) (string, string, error) {
	tail := []string{}
	fullPath := strings.Split(globExpression, string(os.PathSeparator))
	head := fullPath[:]
	for i, slug := range fullPath {
		if isMagic(slug) {
			head = fullPath[:i]
			tail = fullPath[i:]
			break
		}
	}
	return strings.Join(head, string(os.PathSeparator)), strings.Join(tail, string(os.PathSeparator)), nil
}

func (f *folder) GlobTree(g Globber) FileTree {
	tree := FileTree{}
	subFolders := f.Children()

	fyles, _ := justFiles(f.ReadDir("."))

	//	regular files
	for _, fyle := range fyles {
		if g.Match(fyle.(File).FullPath()) {
			tree[fyle] = nil
		}
	}
	//	sub [Folders] (branches)
	for _, subFolder := range subFolders {
		if g.Match(subFolder.FullPath()) {
			tree[subFolder] = subFolder.GlobTree(g)
		} else {
			subTree := subFolder.GlobTree(g)
			//	some folders that don't match the glob pattern must be included
			//	if they have children that match the glob pattern
			if !subTree.IsEmpty() {
				tree[subFolder] = subTree
			}
		}
	}
	return tree
}
