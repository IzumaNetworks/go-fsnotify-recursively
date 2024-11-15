package gorph

import (
	"fmt"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
)

type GorphOp uint8

const (
	UndefinedOp GorphOp = iota
	FsNotifyEvent
	FolderAdded
	FolderRemoved
	FolderRenamed
	FolderMoved
	FolderModified
	FolderUnknown
)

func (gop GorphOp) String() string {
	return []string{"UndefinedOp", "FsNotifyEvent", "FolderAdded", "FolderRemoved", "FolderRenamed", "FolderMoved", "FolderModified", "FolderUnknown"}[gop]
}

type GorphEvent struct {
	NotifyEvent *fsnotify.Event
	Op          GorphOp
	Path        string
	Matches     bool
}

func (gevent GorphEvent) toSSE(namespace string) string {
	id := time.Now().Nanosecond()
	return fmt.Sprintf("id: %x\nevent: %s\ndata: %s\ndata: %s\ndata: %s\ndata: %s\n\n", id, namespace, gevent.Op.String(), gevent.NotifyEvent.Op.String(), gevent.NotifyEvent.Name, gevent.Path)
}

func (gevent GorphEvent) String() string {
	return gevent.toSSE("fs")
}

func NotifyToGorphEvent(g *gorph, fevent *fsnotify.Event) GorphEvent {

	gop := UndefinedOp
	shortPath := g.shortPath(fevent.Name)
	matches, err := doublestar.PathMatch(g.Pattern(), shortPath)
	if err != nil {
		g.Watcher.Errors <- err
	}

	//	wasDir indicates a folder that was being tracked, but was removed from the tree
	wasDir := func(longPath string) bool {
		return g.knownFolders[longPath]
	}

	if IsDir(g.backer, shortPath) {
		switch fevent.Op {
		case fsnotify.Create:
			gop = FolderAdded
		case fsnotify.Remove:
			// @note: this should not happen, because events are fired after the file is removed
			gop = FolderRemoved
		case fsnotify.Rename:
			gop = FolderRenamed
		default:
			gop = FolderUnknown
		}
	} else if wasDir(shortPath) {
		switch fevent.Op {
		case fsnotify.Remove:
			gop = FolderRemoved
		case fsnotify.Rename:
			gop = FolderRenamed
		default:
			gop = FsNotifyEvent
		}
	} else {
		gop = FsNotifyEvent
	}

	if gop == UndefinedOp {
		panic("invalid GorphEvent")
	}
	return GorphEvent{
		NotifyEvent: fevent,
		Op:          gop,
		Path:        shortPath,
		Matches:     matches,
	}
}
