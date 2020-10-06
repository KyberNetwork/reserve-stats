package withdrawalstorage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateWithdrawHistory([]binance.WithdrawHistory, string) error
	GetWithdrawHistory(fromTime, toTime time.Time) (map[string][]binance.WithdrawHistory, error)
	GetLastStoredTimestamp(account string) (time.Time, error)
}
