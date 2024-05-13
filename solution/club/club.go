package club

import (
	"fmt"
	"slices"
	"time"

	"github.com/Razzle131/ComputerClub/solution/client"
	"github.com/Razzle131/ComputerClub/solution/queue"
	"github.com/Razzle131/ComputerClub/solution/table"
)

type Club struct {
	clubTables  []table.Table // clients at tables
	clubClients []string      // all club clients
	clubQueue   queue.Queue   // a queue of clients waiting for a free table
}

// creates new instance of club
func New(tableCount int) Club {
	return Club{make([]table.Table, tableCount), []string{}, queue.Queue{}}
}

// returns index of the table searched by client name,
// set clientName = "" to search free table,
// if there are no such table returns -1
func (c Club) FindTableByClient(clientName string) int {
	for i, t := range c.clubTables {
		if t.TableClient.ClientName == clientName {
			return i
		}
	}
	return -1
}

func (c *Club) AddClient(clientName string) {
	if !slices.Contains(c.clubClients, clientName) {
		c.clubClients = append(c.clubClients, clientName)
	}
}

// removes client name from all clients slice
func (c *Club) DeleteClient(clientName string) {
	if ind := slices.Index(c.clubClients, clientName); ind >= 0 {
		c.clubClients = slices.Delete(c.clubClients, ind, ind+1)
	}
	if ind := slices.Index(c.clubQueue, clientName); ind >= 0 {
		c.clubQueue = slices.Delete(c.clubQueue, ind, ind+1)
	}
}

// checks if client is already came to club
func (c Club) IsInClub(clientName string) bool {
	return slices.Contains(c.clubClients, clientName)
}

// checks if given table is busy. If yes, returns true; if not, returns false
func (c Club) IsTableBusy(tableId int) bool {
	if tableId >= len(c.clubTables) {
		return true
	}
	return c.clubTables[tableId].TableClient.ClientName != ""
}

// ends previous client session (if exists) and start new
func (c *Club) StartNewSession(curEventClient client.Client, price int, tableId int) {
	if !c.IsInClub(curEventClient.ClientName) {
		return
	}

	if leavedClientTable := c.FindTableByClient(curEventClient.ClientName); leavedClientTable >= 0 {
		c.clubTables[leavedClientTable].EndTableClientSession(curEventClient.StartedTime, price)
	}
	c.clubTables[tableId].TableClient = curEventClient
}

// checks if queue is over limit size given in file
func (c Club) QueueIsTooBig(maxSize int) bool {
	return len(c.clubQueue) > maxSize
}

// adds client to queue
func (c *Club) AddClientToQueue(clientName string) {
	if !slices.Contains(c.clubClients, clientName) {
		c.clubClients = append(c.clubClients, clientName)
	}
	if !slices.Contains(c.clubQueue, clientName) {
		c.clubQueue.Push(clientName)
	}
}

// ends given client session and starts new if queue is not empty
func (c *Club) EndClientSession(curClient client.Client, price int) {
	c.DeleteClient(curClient.ClientName)

	if leavedClientTableInd := c.FindTableByClient(curClient.ClientName); leavedClientTableInd >= 0 {
		leavedTable := &c.clubTables[leavedClientTableInd]
		leavedTable.EndTableClientSession(curClient.StartedTime, price)
		if len(c.clubQueue) > 0 {
			newClientName := c.clubQueue.Pop()
			fmt.Println(curClient.StartedTime.Format("15:04"), 12, newClientName, leavedClientTableInd+1)

			leavedTable.TableClient = client.Client{ClientName: newClientName, StartedTime: curClient.StartedTime}
		}
	}
}

// inits closure proccess in club
func (c *Club) Close(endTime time.Time, price int, inputTimeFormat string) {
	slices.Sort(c.clubClients)
	for _, clientName := range c.clubClients {
		fmt.Println(endTime.Format(inputTimeFormat), 11, clientName)
		if leavedClientTable := c.FindTableByClient(clientName); leavedClientTable >= 0 {
			leavedTable := &c.clubTables[leavedClientTable]
			leavedTable.EndTableClientSession(endTime, price)
		}
	}
	c.clubClients = []string{}
	c.clubQueue = queue.Queue{}
}

// returns formatted income for all tables by slice
func (c Club) GetIncome() []string {
	var res []string
	for id, table := range c.clubTables {
		tableHours := int(table.TableTotalTime.Hours())
		tableMinutes := int(table.TableTotalTime.Minutes()) % 60
		res = append(res, fmt.Sprint(id+1)+" "+fmt.Sprint(table.TableTotalIncome)+" "+fmt.Sprintf("%02d:%02d", tableHours, tableMinutes))
	}
	return res
}
