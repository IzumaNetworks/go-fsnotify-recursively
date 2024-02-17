package fsnotifyr

import (
	"cmp"
	"io/fs"
	"slices"
	"strings"

	"github.com/xlab/treeprint"
)

// FileTree is a tree containing files and directories
// files and empty directories are leaf nodes
type FileTree map[fs.DirEntry]FileTree

type record struct {
	dirEntry fs.DirEntry
	branch   FileTree
}

// entries returns key-value pairs as sorted slice
func (t FileTree) entries() []record {
	kvps := []record{}
	for k, v := range t {
		kvps = append(kvps, record{k, v})
	}
	//	alpha sort, by file name
	slices.SortFunc(kvps, func(a, b record) int {
		return cmp.Compare(a.dirEntry.Name(), b.dirEntry.Name())
	})
	return kvps
}

func (t FileTree) String() string {
	rootNode := treeprint.New()
	//	recursive function
	var fn func(treeprint.Tree, FileTree)
	fn = func(node treeprint.Tree, tree FileTree) {
		for _, entry := range tree.entries() {
			k := entry.dirEntry
			v := entry.branch
			if k.IsDir() {
				subNode := node.AddBranch(k.Name())
				fn(subNode, v)
			} else {
				if k.Type().IsRegular() {
					node.AddNode(k.Name())
				}
			}
		}
	}
	fn(rootNode, t)
	return strings.TrimSpace(rootNode.String())
}

func (f *folder) FileTree(includeFiles bool) FileTree {
	tree := FileTree{}
	subFolders := f.Children()
	if includeFiles {
		fyles, _ := justFiles(f.ReadDir("."))
		//	regular files
		for _, fyle := range fyles {
			tree[fyle] = nil
		}
	}
	//	sub [Folders] (branches)
	for _, subFolder := range subFolders {
		tree[subFolder] = subFolder.FileTree(includeFiles)
	}
	return tree
}
