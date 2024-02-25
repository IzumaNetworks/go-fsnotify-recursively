package fsnotifyr

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type FolderEvent struct {
	Prefix        string
	File          File
	Folder        Folder
	Op            string
	FsNotifyEvent *fsnotify.Event
}

func (fe FolderEvent) String() string {
	m := map[string]string{
		"prefix": fe.Prefix,
		"path":   fe.Folder.FullPath(),
		"op":     fe.Op,
	}
	j, _ := json.Marshal(m)
	return string(j)
}

type WatchTree struct {
	rootFolder   Folder
	prefix       string
	globPattern  Globber
	Watcher      *fsnotify.Watcher
	FolderEvents chan FolderEvent
}

func (w *WatchTree) Filesystem() fs.FS {
	return w.RootFolder().Filesystem()
}

func (w *WatchTree) RootFolder() Folder {
	return w.rootFolder
}

func (w *WatchTree) Globber() Globber {
	return w.globPattern
}

func (w *WatchTree) AddFolder(f Folder) {
	//	prepend the prefix and pass the event to underlying fsnotify
	fullyFullPath := filepath.Join(w.prefix, f.FullPath())
	fmt.Println(fullyFullPath)
	w.Watcher.Add(fullyFullPath)
	ev := FolderEvent{w.prefix, f, "add", nil}

	//	broadcast self
	w.FolderEvents <- ev

	//	add children
	for _, child := range f.Children() {
		w.AddFolder(child)
	}
}

func (w *WatchTree) RemoveFolder(f Folder) {
	ev := FolderEvent{w.prefix, f, "destroy", nil}
	w.FolderEvents <- ev
	f.Destroy()
}

func StringToDirentry(path string, prefix string, filesystem fs.FS) (fs.DirEntry, error) {
	fullPath := filepath.Join(prefix, path)
	f, err := filesystem.Open(fullPath)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return fs.FileInfoToDirEntry(stat), nil
}

func (w *WatchTree) Listen() (chan FolderEvent, chan fsnotify.Event, chan error) {

	//	add all folders
	go func() {
		w.AddFolder(w.RootFolder())
		for ev := range w.Watcher.Events {
			dirEntry, err := StringToDirentry(ev.Name, w.prefix, w.Filesystem())
			if err != nil {
				panic(err)
			}
			if dirEntry.IsDir() {
				//	create a new folder
				//	but first we need to locate the proper place in the tree
				newFolder, err := AddDescendantFolder(ev.Name, w.RootFolder())
				if err != nil {
					panic(err)
				}
				w.AddFolder(newFolder)
			} else {
				//	this is a file, so wrap and pass
				fmt.Println("@todo")
			}
		}
	}()

	return w.FolderEvents, w.Watcher.Events, w.Watcher.Errors
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

	folderEventsChannel := make(chan FolderEvent)

	w := WatchTree{
		prefix:       pathPrefix,
		globPattern:  compiledGlob,
		rootFolder:   rootFolder,
		Watcher:      mainWatcher,
		FolderEvents: folderEventsChannel,
	}

	return &w, nil

}
