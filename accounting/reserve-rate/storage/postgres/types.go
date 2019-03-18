package postgres

import (
	"time"
)

type ratesFromDB struct {
	Time    time.Time `db:"time"`
	Token   string    `db:"token"`
	Base    string    `db:"base"`
	Rate    float64   `db:"rate"`
	Reserve string    `db:"reserve"`
}

type usdRateFromDB struct {
	Time time.Time `db:"time"`
	Rate float64   `db:"rate"`
}
