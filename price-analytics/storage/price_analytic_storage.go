package storage

import "github.com/go-pg/pg"

//PriceAnalyticDB represent db for price analytic data
type PriceAnalyticDB struct {
	db *pg.DB
}

//StorePriceAnalytic store information from send from analytic to stat
func (pad *PriceAnalyticDB) StorePriceAnalytic() error {
	return nil
}

//NewPriceStorage return new storage for price analytics
func NewPriceStorage() {

}
