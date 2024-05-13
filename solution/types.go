package Solution

import (
	"fmt"
	"math"
	"slices"
	"time"
)

type client struct {
	clientName  string
	startedTime time.Time // time when the client starts his session
}

func (c client) getPayment(endSessionTime time.Time, price int) int {
	var durInHours float64 = endSessionTime.Sub(c.startedTime).Hours()
	if durInHours < 0 {
		return 0
	}
	return int(math.Ceil(durInHours)) * price
}

type table struct {
	tableTotalTime   time.Duration
	tableTotalIncome int
	tableClient      client
}

// accumulates table time and income and resets user of the table
func (t *table) endClientSession(endSessionTime time.Time, price int) {
	var dur time.Duration = endSessionTime.Sub(t.tableClient.startedTime)
	if dur > 0 {
		t.tableTotalTime += dur
	}

	t.tableTotalIncome += t.tableClient.getPayment(endSessionTime, price)
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

// removes client name from all clients slice
func (c *club) deleteClient(clientName string) {
	if ind := slices.Index(c.cameClients, clientName); ind >= 0 {
		c.cameClients = slices.Delete(c.cameClients, ind, ind+1)
	}
}

// inits closure proccess in club, endTimeStr needs if we have endTime input = "24:00"
func (c *club) close(endTimeStr string, endTime time.Time, price int) {
	slices.Sort(c.cameClients)
	for _, client := range c.cameClients {
		fmt.Println(endTimeStr, 11, client)
		if leavedClientTable := c.findTableByClient(client); leavedClientTable >= 0 {
			leavedTable := &c.clubClients[leavedClientTable]
			leavedTable.endClientSession(endTime, price)
		}
	}
	c.cameClients = []string{}
	c.clubQueue = queue{}
}

// checks if previous event was in new day or its time is bigger than new event, so modifies current date by delta days
func shiftDate(prevEventDate time.Time, curEventDate *time.Time) {
	if prevEventDate != (time.Time{}) {
		if delta := prevEventDate.Day() - curEventDate.Day(); delta > 0 {
			*curEventDate = curEventDate.AddDate(0, 0, delta)
		}
		if prevEventDate.Compare(*curEventDate) > 0 {
			*curEventDate = curEventDate.AddDate(0, 0, 1)
		}
	}
}
