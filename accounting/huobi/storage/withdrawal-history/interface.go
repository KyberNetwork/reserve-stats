package withdrawalhistory

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

//Interface is the storage for withdrawal history
type Interface interface {
	UpdateWithdrawHistory(withdraws []huobi.WithdrawHistory) error
	GetWithdrawHistory(from, to time.Time) ([]huobi.WithdrawHistory, error)
	GetLastIDStored() (uint64, error)
}
