package Solution

import (
	"math"
	"slices"
	"strings"
	"time"
)

type client struct {
	clientName  string
	startedTime time.Time // time when the client starts his session
}

func (c client) getPayment(endSessionTime time.Time) int {
	return int(math.Ceil(endSessionTime.Sub(c.startedTime).Hours())) * price
}

type table struct {
	tableTotalTime   time.Duration
	tableTotalIncome int
	tableClient      client
}

// accumulates table time and income and resets user of the table
func (t *table) endClientSession(endSessionTime time.Time) {
	t.tableTotalTime += endSessionTime.Sub(t.tableClient.startedTime)
	t.tableTotalIncome += t.tableClient.getPayment(endSessionTime)
	t.tableClient = client{}
}

type queue []string

// returns first element of queue and deletes it from queue
func (q *queue) pop() string {
	var res string = (*q)[0]
	*q = (*q)[1:]
	return res
}

// adds value to end of queue
func (q *queue) push(val string) {
	*q = append(*q, val)
}

type club struct {
	clubClients []table  // clients at tables
	cameClients []string // all club clients
	clubQueue   queue    // a queue of clients waiting for a free table
}

// returns index of the table searched by client name,
// set clientName = "" to search free table,
// if there are no such table returns -1
func (c club) findTableByClient(clientName string) int {
	for i, t := range c.clubClients {
		if t.tableClient.clientName == clientName {
			return i
		}
	}
	return -1
}

// removes client name from clients slice
func (c *club) deleteClient(clientName string) {
	if ind := slices.Index(c.cameClients, clientName); ind >= 0 {
		c.cameClients = slices.Delete(c.cameClients, ind, ind+1)
	}
}

// checks if there are events before that in standart would be later than checked event, this makes this happen next day
func checkIfThisIsNewDay(i int, input []string, curEventTime time.Time) bool {
	for j := i; j > 2; j-- {
		prevEventTime, err := time.Parse(timeLayout, strings.Split(input[j], " ")[0])
		if err == nil && prevEventTime.Compare(curEventTime) > 0 {
			return true
		}
	}
	return false
}
