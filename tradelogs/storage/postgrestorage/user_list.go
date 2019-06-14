package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

//TODO implement this
func (tldb *TradeLogDB) GetUserList(fromTime, toTime time.Time, timezone int8) ([]common.UserInfo, error) {
	return nil, nil
}
