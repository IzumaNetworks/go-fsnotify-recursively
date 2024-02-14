package rwatch

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

type WatchTree struct {
	rootFolder  Folder
	prefix      string
	globPattern Globber
	Watcher     *fsnotify.Watcher
}

func NewWatchTree(rootGlob string) (*WatchTree, error) {

	pathPrefix, shortGlob, err := GlobParent(rootGlob)
	if err != nil {
		return nil, err
	}

	fileSystem := os.DirFS(pathPrefix)
	compiledGlob := NewGlobber(shortGlob)

	mainWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rootFolder, err := NewFolder(fileSystem, ".", nil)
	if err != nil {
		return nil, err
	}

	w := WatchTree{
		prefix:      pathPrefix,
		globPattern: compiledGlob,
		rootFolder:  rootFolder,
		Watcher:     mainWatcher,
	}

	return &w, nil

}
