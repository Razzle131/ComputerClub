package Solution

import (
	"testing"
	"time"
)

func TestGetPayment(t *testing.T) {
	var c client = client{"", time.Time{}}

	if c.getPayment(time.Time{}, 1) != 0 {
		t.Fatal("payment with zero duration must be zero")
	}

	c = client{"", time.Time{}.AddDate(0, 0, 1)}
	if c.getPayment(time.Time{}, 1) < 0 {
		t.Fatal("payment cant be negative")
	}

	c = client{"", time.Time{}.AddDate(0, 0, -1)}
	if c.getPayment(time.Time{}, 1) != 24 {
		t.Fatal("payment calculation gone wrong with int number of hours")
	}

	c = client{"", time.Time{}.Add(time.Minute * -25)}
	if c.getPayment(time.Time{}, 1) != 1 {
		t.Fatal("payment calculation gone wrong with float number of hours over half")
	}

	c = client{"", time.Time{}.Add(time.Minute * -35)}
	if c.getPayment(time.Time{}, 1) != 1 {
		t.Fatal("payment calculation gone wrong with float number of hours less than half")
	}
}

func TestEndClientSession(t *testing.T) {
	var tab table = table{}

	tab.endClientSession(time.Time{}.AddDate(0, 0, -1), 1)
	if tab.tableTotalIncome < 0 || tab.tableTotalTime < 0 {
		t.Fatal("table income or table time is decreased")
	}

	tab = table{time.Hour * 0, 0, client{"testclient", time.Time{}}}
	tab.endClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	if tab.tableTotalTime != time.Hour*24 {
		t.Fatal("table time is not calculated rigth")
	}
	if tab.tableTotalIncome != 24 {
		t.Fatal("table income is not calculated rigth")
	}
	if tab.tableClient != (client{}) {
		t.Fatal("client is not reseted")
	}

	tab = table{}
	tab.endClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	tab.endClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	if tab.tableTotalTime != time.Hour*48 {
		t.Fatal("table time is not accumulated")
	}
	if tab.tableTotalIncome != 48 {
		t.Fatal("table income is not accumulated")
	}
}

func TestFindTableByClient(t *testing.T) {
	var c club = club{make([]table, 1), []string{}, queue{}}

	if c.findTableByClient("") != 0 {
		t.Fatal("cant find free table while there are some")
	}

	if c.findTableByClient("someone?") >= 0 {
		t.Fatal("found table with not existing client")
	}

	c.clubClients[0] = table{time.Hour * 0, 0, client{"client", time.Time{}}}
	if c.findTableByClient("client") != 0 {
		t.Fatal("not found table with existing client")
	}

	if c.findTableByClient("") >= 0 {
		t.Fatal("found free table while there are no such")
	}

	if c.findTableByClient("someone?") >= 0 {
		t.Fatal("found table with not existing client")
	}
}

func TestDeleteClient(t *testing.T) {
	var c club = club{make([]table, 1), []string{}, queue{}}

	c.deleteClient("someone?") // doesnt respond with panic if cameClients is empty

	c.cameClients = append(c.cameClients, "client")
	c.deleteClient("someone?") // doesnt respond with panic if client is not in slice
	c.deleteClient("client")

	if len(c.cameClients) != 0 {
		t.Fatal("does not delete existing client")
	}
}

func TestShiftDate(t *testing.T) {
	events := []string{"18:00", "01:00", "00:30"}
	curEventTime, _ := time.Parse(timeLayout, events[2])
	curEventTimeCopy := curEventTime
	prevEventTime, _ := time.Parse(timeLayout, events[1])

	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy.AddDate(0, 0, 1) {
		t.Fatal("no time shift occured for 1 day interval")
	}

	curEventTime, _ = time.Parse(timeLayout, events[2])
	curEventTimeCopy = curEventTime
	prevEventTime, _ = time.Parse(timeLayout, events[1])
	firstEventTime, _ := time.Parse(timeLayout, events[0])

	shiftDate(firstEventTime, &prevEventTime)
	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy.AddDate(0, 0, 2) {
		t.Fatal("no time shift occured for 2 days interval")
	}

	events = []string{"18:00", "01:00", "00:30"}
	curEventTime, _ = time.Parse(timeLayout, events[0])
	curEventTimeCopy = curEventTime
	prevEventTime, _ = time.Parse(timeLayout, events[1])

	shiftDate(prevEventTime, &curEventTime)
	if curEventTime != curEventTimeCopy {
		t.Fatal("time shift occured while it should not be")
	}
}
