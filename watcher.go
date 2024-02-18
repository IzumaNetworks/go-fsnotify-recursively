package fsnotifyr

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

func (w *WatchTree) RootFolder() Folder {
	return w.rootFolder
}

func (w *WatchTree) Globber() Globber {
	return w.globPattern
}

func NewWatchTree(rootGlob string) (*WatchTree, error) {

	compiledGlob, err := NewGlobber(rootGlob)
	pathPrefix := compiledGlob.Root()
	fileSystem := os.DirFS(pathPrefix)

	if err != nil {
		return nil, err
	}

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
