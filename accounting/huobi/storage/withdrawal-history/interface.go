package withdrawalhistory

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

//Interface is the storage for withdrawal history
type Interface interface {
	UpdateWithdrawHistory(withdraws []huobi.WithdrawHistory, account string) error
	GetWithdrawHistory(from, to time.Time) (map[string][]huobi.WithdrawHistory, error)
	GetLastIDStored(account string) (uint64, error)
}
