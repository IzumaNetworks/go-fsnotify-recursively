package rwatch

import (
	"io/fs"
	"strings"

	"github.com/xlab/treeprint"
)

// FileTree is a tree containing files and directories
// files and empty directories are leaf nodes
type FileTree map[fs.DirEntry]FileTree

func (t FileTree) String() string {
	rootNode := treeprint.New()

	//	recursive function
	var fn func(treeprint.Tree, FileTree)
	fn = func(node treeprint.Tree, tree FileTree) {
		for k, v := range tree {
			if k.IsDir() {
				subNode := node.AddBranch(k.Name())
				fn(subNode, v)
			} else {
				//	omit symlinks, FIFO pipes, device files, etc
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
	fyles, _ := justFiles(f.ReadDir("."))

	//	regular files
	for _, fyle := range fyles {
		tree[fyle] = nil
	}

	//	sub [Folders] (branches)
	for _, subFolder := range subFolders {
		tree[subFolder] = subFolder.FileTree(includeFiles)
	}

	return tree
}
