package main

import (
	"flag"
	"fmt"

	gorph "github.com/sean9999/go-fsnotify-recursively"
)

func main() {

	globExpressionPtr := flag.String("glob", "**", "the pattern to watch")

	flag.Parse()

	g, err := gorph.New(*globExpressionPtr)
	if err != nil {
		panic(err)
	}

	evs, ers := g.Listen()

	for {
		select {
		case ev, ok := <-evs:
			if !ok {
				sseEvent := createSSE("gorph/exit", "exiting")
				fmt.Println(sseEvent)
				return
			}

			//	omit events for regular files where glob does not match
			if ev.Op != gorph.FsNotifyEvent || ev.Matches {
				fmt.Println(gorphAsSSE(ev))
			}

		case er, ok := <-ers:
			if ok {
				fmt.Println(createSSE("gorph/error", er.Error()))
			} else {
				return
			}
		}
	}

}
