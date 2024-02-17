package fsnotifyr

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

// helper function to clean up fs.ReadDir()
// also sorts, to assist in predictability
func justFiles(entries []fs.DirEntry, err error) ([]fs.DirEntry, error) {
	fyles := []fs.DirEntry{}
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			fyles = append(fyles, entry)
		}
	}
	//	alpha sort
	// slices.SortStableFunc(fyles, func(a, b fs.DirEntry) int {
	// 	return cmp.Compare(a.Name(), b.Name())
	// })
	return fyles, nil
}

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

func NewWatcherTree(rootGlob string) map[string]bool {

	r := map[string]bool{}

	root, rest, err := ComponentizeGlobString(rootGlob)
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
