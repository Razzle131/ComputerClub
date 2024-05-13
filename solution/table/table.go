package table

import (
	"time"

	client "github.com/Razzle131/ComputerClub/solution/client"
)

type Table struct {
	TableTotalTime   time.Duration
	TableTotalIncome int
	TableClient      client.Client
}

// accumulates table time and income and resets user of the table
func (t *Table) EndTableClientSession(endSessionTime time.Time, price int) {
	var dur time.Duration = endSessionTime.Sub(t.TableClient.StartedTime)
	if dur > 0 {
		t.TableTotalTime += dur
	}

	t.TableTotalIncome += t.TableClient.GetPayment(endSessionTime, price)
	t.TableClient = client.Client{}
}
