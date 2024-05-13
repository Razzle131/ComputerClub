package event

import (
	"fmt"
	"strings"
)

type Event struct {
	EventTime       string
	EventId         string
	EventClient     string
	EventTable      string
	EventErrMessage string
}

func (e Event) PrintEvent() {
	if e.EventErrMessage != "" {
		e.EventId = "13"
		e.EventClient = ""
		e.EventTable = ""
	}
	eventMessage := fmt.Sprintf("%v %v %v %v %v", e.EventTime, e.EventId, e.EventClient, e.EventTable, e.EventErrMessage)

	for _, val := range strings.Split(eventMessage, " ") {
		if val != "" {
			fmt.Printf("%v ", val)
		}
	}
	fmt.Println()
}
