package main

import (
	"fmt"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

func main() {

	watcher, err := fsnotifyr.NewWatchTree("testdata/**")
	if err != nil {
		panic(err)
	}

	lifecycleEvents, fsEvents, errors := watcher.Listen()

	for {
		select {
		case ev, ok := <-fsEvents:
			fmt.Println("fs", ev, ok)
		case err, ok := <-errors:
			fmt.Println("err", err, ok)
		case life, ok := <-lifecycleEvents:
			fmt.Println(life, ok)
		}
	}

}
