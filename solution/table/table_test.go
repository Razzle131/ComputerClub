package table

import (
	"testing"
	"time"

	"github.com/Razzle131/ComputerClub/solution/client"
)

func TestEndClientSession(t *testing.T) {
	var tab Table = Table{}

	tab.EndTableClientSession(time.Time{}.AddDate(0, 0, -1), 1)
	if tab.TableTotalIncome < 0 || tab.TableTotalTime < 0 {
		t.Fatal("table income or table time is decreased")
	}

	tab = Table{time.Hour * 0, 0, client.Client{ClientName: "testclient", StartedTime: time.Time{}}}
	tab.EndTableClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	if tab.TableTotalTime != time.Hour*24 {
		t.Fatal("table time is not calculated rigth")
	}
	if tab.TableTotalIncome != 24 {
		t.Fatal("table income is not calculated rigth")
	}
	if tab.TableClient != (client.Client{}) {
		t.Fatal("client is not reseted")
	}

	tab = Table{}
	tab.EndTableClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	tab.EndTableClientSession(time.Time{}.AddDate(0, 0, 1), 1)
	if tab.TableTotalTime != time.Hour*48 {
		t.Fatal("table time is not accumulated")
	}
	if tab.TableTotalIncome != 48 {
		t.Fatal("table income is not accumulated")
	}
}
