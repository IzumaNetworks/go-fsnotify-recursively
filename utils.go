package rwatch

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gobwas/glob"
)

type Jimenju struct {
	MotherWatcher fsnotify.Watcher
	RootGlob      glob.Glob
	Prefix        string
}

// }

// if this slug (portion of a path) is a glob pattern, it's magic
func isMagic(slug string) bool {
	index := -1
	magical_runes := []rune{'*', '?', '[', '{'}
	for _, char := range magical_runes {
		if index = strings.IndexRune(slug, char); index > -1 {
			return true
		}
	}
	return false
}

func GlobParent(globExpression string) (string, string, error) {
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

func NewWatcherTree(rootGlob string) map[string]bool {

	r := map[string]bool{}

	root, rest, err := GlobParent(rootGlob)
	if err != nil {
		panic(err)
	}

	fileSystem := os.DirFS(root)
	g := glob.MustCompile(rest, os.PathSeparator)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fullPath, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		matches := g.Match(fullPath)
		//fmt.Println(relativePath, matches)
		r[fullPath] = matches
		return nil
	})

	return r

}
