package withdrawstorage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateWithdrawHistory(map[string]binance.WithdrawHistory) error
	GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error)
}
