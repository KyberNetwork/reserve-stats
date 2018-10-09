package stats

import (
	"errors"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/ethrate"
	"github.com/KyberNetwork/reserve-stats/lib/utils"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/go-pg/pg"
)

//UserStats represent stats for an user
type UserStats struct {
	cmcEthUSDRate ethrate.EthUSDRate
	userStorage   *storage.UserDB
}

//GetTxCapByAddress return user Tx limit by wei
//return true if address kyced, and return false if address is non-kyced
func (us UserStats) GetTxCapByAddress(addr string) (*big.Int, bool, error) {
	_, err := us.userStorage.GetUserInfo(addr)
	var usdCap float64
	kyced := true
	usdCap = common.KycedCap().DailyLimit
	if err != nil {
		if err == pg.ErrNoRows {
			usdCap = common.NonKycedCap().TxLimit
			kyced = false
		} else {
			return nil, false, err
		}
	}
	timepoint := time.Now()
	rate := us.cmcEthUSDRate.GetUSDRate(timepoint)
	var txLimit *big.Int
	if rate == 0 {
		return txLimit, kyced, errors.New("cannot get eth usd rate from cmc")
	}
	ethLimit := usdCap / rate
	txLimit = utils.EthToWei(ethLimit)
	return txLimit, kyced, nil
}

//StoreUserInfo store user info pushed from dashboard
func (us UserStats) StoreUserInfo(userData common.UserData) error {
	return us.userStorage.StoreUserInfo(userData)
}

//NewUserStats return new user stats instance
func NewUserStats(cmc ethrate.EthUSDRate, storage *storage.UserDB) *UserStats {
	return &UserStats{
		cmcEthUSDRate: cmc,
		userStorage:   storage,
	}
}
