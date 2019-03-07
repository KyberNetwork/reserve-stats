package postgres

import (
	"time"
)

type ratesFromDB struct {
	Time    time.Time `db:"time"`
	Token   string    `db:"token"`
	Base    string    `db:"base"`
	BuyRate float64   `db:"buy_rate"`
	Reserve string    `db:"reserve"`
}
