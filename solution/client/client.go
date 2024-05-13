package client

import (
	"math"
	"time"
)

type Client struct {
	ClientName  string
	StartedTime time.Time // time when the client starts his session
}

func New(clientName string, startedTime time.Time) Client {
	return Client{ClientName: clientName, StartedTime: startedTime}
}

func (c Client) GetPayment(endSessionTime time.Time, price int) int {
	var durInHours float64 = endSessionTime.Sub(c.StartedTime).Hours()
	if durInHours < 0 {
		return 0
	}
	return int(math.Ceil(durInHours)) * price
}
