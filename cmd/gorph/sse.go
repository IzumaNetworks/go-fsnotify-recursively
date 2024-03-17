package main

import (
	"fmt"
	"strings"
	"time"

	gorph "github.com/sean9999/go-fsnotify-recursively"
)

type ServerSentEvent struct {
	tz        int
	eventType string
	data      []string
}

func (sse ServerSentEvent) String() string {
	lines := []string{fmt.Sprintf("id:\t%x", sse.tz)}

	if len(sse.eventType) > 1 {
		lines = append(lines, fmt.Sprintf("event:\t%s", sse.eventType))
	}

	for _, datum := range sse.data {
		lines = append(lines, fmt.Sprintf("data:\t%s", datum))
	}

	//	SSE ends with two line breaks
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func createSSE(eventType string, datas ...string) ServerSentEvent {
	r := ServerSentEvent{
		tz:        time.Now().Nanosecond(),
		eventType: eventType,
		data:      datas,
	}
	return r
}

func gorphAsSSE(gevent gorph.GorphEvent) ServerSentEvent {
	eventType := fmt.Sprintf("gorph/%s", gevent.Op)
	matches := "GlobMatch:\tNO"
	if gevent.Matches {
		matches = "GlobMatch:\tYES"
	}
	r := createSSE(eventType, gevent.Path, matches, gevent.NotifyEvent.String())
	return r
}
