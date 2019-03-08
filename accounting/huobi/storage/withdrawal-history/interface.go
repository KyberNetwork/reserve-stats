package withdrawalhistory

import "github.com/KyberNetwork/reserve-stats/lib/huobi"

//Interface is the storage for withdrawal history
type Interface interface {
	UpdateWithdrawHistory(withdraw huobi.WithdrawHistory) error
}
