package main

import (
	"fmt"
	"os"

	gorph "go.izuma.io/go-fsnotify-recursively"
)

func main() {

	globExpression := "**"
	if len(os.Args) > 1 {
		globExpression = os.Args[1]
	}

	g, err := gorph.New(globExpression)
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

			//	omit events for regular files where glob doesn't match
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
