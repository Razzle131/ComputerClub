package client

import (
	"testing"
	"time"
)

func TestGetPayment(t *testing.T) {
	var c Client = Client{"", time.Time{}}

	if c.GetPayment(time.Time{}, 1) != 0 {
		t.Fatal("payment with zero duration must be zero")
	}

	c = Client{"", time.Time{}.AddDate(0, 0, 1)}
	if c.GetPayment(time.Time{}, 1) < 0 {
		t.Fatal("payment cant be negative")
	}

	c = Client{"", time.Time{}.AddDate(0, 0, -1)}
	if c.GetPayment(time.Time{}, 1) != 24 {
		t.Fatal("payment calculation gone wrong with int number of hours")
	}

	c = Client{"", time.Time{}.Add(time.Minute * -25)}
	if c.GetPayment(time.Time{}, 1) != 1 {
		t.Fatal("payment calculation gone wrong with float number of hours over half")
	}

	c = Client{"", time.Time{}.Add(time.Minute * -35)}
	if c.GetPayment(time.Time{}, 1) != 1 {
		t.Fatal("payment calculation gone wrong with float number of hours less than half")
	}

	c = Client{"", time.Time{}.Add(time.Minute * -30)}
	if c.GetPayment(time.Time{}, 1) != 1 {
		t.Fatal("payment calculation gone wrong with float number of hours less than half")
	}
}
