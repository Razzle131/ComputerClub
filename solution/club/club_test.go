package club

import (
	"testing"
	"time"

	"github.com/Razzle131/ComputerClub/solution/client"
	"github.com/Razzle131/ComputerClub/solution/queue"
	"github.com/Razzle131/ComputerClub/solution/table"
)

func TestNew(t *testing.T) {
	var c Club = New(1)

	if len(c.clubQueue) != 0 {
		t.Fatal("queue created wrong")
	}

	if len(c.clubClients) != 0 {
		t.Fatal("clients slice created wrong")
	}

	if c.clubTables[0].TableTotalIncome != 0 || c.clubTables[0].TableTotalTime != time.Hour*0 || c.clubTables[0].TableClient != (client.Client{}) {
		t.Fatal("default values of a table created wrong")
	}

	c = New(0)
	if len(c.clubQueue) != 0 {
		t.Fatal("queue created wrong")
	}

	if len(c.clubClients) != 0 {
		t.Fatal("clients slice created wrong")
	}

	if len(c.clubTables) != 0 {
		t.Fatal("tables slice created wrong")
	}
}

func TestFindTableByClient(t *testing.T) {
	var c Club = Club{make([]table.Table, 1), []string{}, queue.Queue{}}

	if c.FindTableByClient("") != 0 {
		t.Fatal("cant find free table while there are some")
	}

	if c.FindTableByClient("someone?") >= 0 {
		t.Fatal("found table with not existing client")
	}

	c.clubTables[0] = table.Table{TableTotalTime: time.Hour * 0, TableTotalIncome: 0, TableClient: client.New("client", time.Time{})}
	if c.FindTableByClient("client") != 0 {
		t.Fatal("not found table with existing client")
	}

	if c.FindTableByClient("") >= 0 {
		t.Fatal("found free table while there are no such")
	}

	if c.FindTableByClient("someone?") >= 0 {
		t.Fatal("found table with not existing client")
	}
}

func TestAddClient(t *testing.T) {
	var c Club = New(0)

	c.AddClient("client")
	if c.clubClients[0] != "client" {
		t.Fatal("not found added client")
	}

	c = New(0)
	c.AddClient("client")
	c.AddClient("client")
	if len(c.clubClients) > 1 {
		t.Fatal("added same client twice")
	}
}

func TestDeleteClient(t *testing.T) {
	var c Club = Club{make([]table.Table, 1), []string{}, queue.Queue{}}

	c.DeleteClient("someone?") // doesnt respond with panic if cameClients is empty

	c.clubClients = append(c.clubClients, "client")
	c.DeleteClient("someone?") // doesnt respond with panic if client is not in slice
	c.DeleteClient("client")

	if len(c.clubClients) != 0 {
		t.Fatal("does not delete existing client")
	}
}

func TestIsTableBusy(t *testing.T) {
	var c Club = Club{make([]table.Table, 1), []string{}, queue.Queue{}}

	if c.IsTableBusy(0) {
		t.Fatal("this table should be empty")
	}

	c.clubTables[0].TableClient.ClientName = "client"
	if !c.IsTableBusy(0) {
		t.Fatal("this table should not be empty")
	}

	if !c.IsTableBusy(10) {
		t.Fatal("table which is over table count, should not be empty")
	}
}

func TestStartNewSession(t *testing.T) {
	c := New(1)

	c.StartNewSession(client.New("client", time.Time{}), 1, 0)
	if c.clubTables[0].TableClient.ClientName != "" {
		t.Fatal("started session of not existing client")
	}

	c = New(1)
	c.AddClient("client")
	c.StartNewSession(client.New("client", time.Time{}), 1, 0)
	if c.clubTables[0].TableClient.ClientName == "" {
		t.Fatal("doesnt started client session")
	}
}

func TestEndClientSession(t *testing.T) {
	c := New(1)

	c.AddClient("client")
	oldClientEvent := client.New("client", time.Time{}.AddDate(0, 0, -1))
	newClientEvent := client.New("client", time.Time{})
	c.StartNewSession(oldClientEvent, 1, 0)
	c.EndClientSession(newClientEvent, 1)

	if len(c.clubClients) != 0 {
		t.Fatal("client gone, but in clients slice he still occures")
	}

	if c.clubTables[0].TableTotalIncome != 24 || c.clubTables[0].TableTotalTime != time.Hour*24 {
		t.Fatal("client ended session, but metrics of the table wasnt updated")
	}

	c = New(1)
	c.EndClientSession(newClientEvent, 1) // doesnt throw error if client doesnt exist

	if c.clubTables[0].TableTotalIncome != 0 || c.clubTables[0].TableTotalTime != time.Hour*0 {
		t.Fatal("counted metrics of not existing client")
	}
}

func TestClose(t *testing.T) {
	c := New(1)
	newClient := client.New("client", time.Time{}.AddDate(-1, 0, 0))
	closeTime, _ := time.Parse("15:04", "19:00")

	c.AddClient("client")
	c.AddClientToQueue("client2")
	c.StartNewSession(newClient, 1, 0)
	c.Close(closeTime, 1, "15:04")

	if len(c.clubClients) > 0 {
		t.Fatal("after closure there are clients in clients slice")
	}

	if len(c.clubQueue) > 0 {
		t.Fatal("after closure there are clients in queue")
	}

	if c.clubTables[0].TableClient != (client.Client{}) || c.clubTables[0].TableTotalIncome != 19 || c.clubTables[0].TableTotalTime != time.Hour*19 {
		t.Fatal("table fields are not setted rigth way")
	}

}
