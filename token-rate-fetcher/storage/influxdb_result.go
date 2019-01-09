package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/common"
	schema "github.com/KyberNetwork/reserve-stats/token-rate-fetcher/storage/schema/tokenrate"
)

//LastTimePoint return first  time point in db or error if it is empty.
func (is *InfluxStorage) LastTimePoint(providerName, tokenID, currencyID string) (time.Time, error) {
	measurementName := fmt.Sprintf("%s_%s", common.GetTokenSymbolFromProviderNameTokenID(providerName, tokenID), currencyID)
	stmt := fmt.Sprintf(`SELECT LAST("%s") FROM "%s"`,
		schema.Rate.String(),
		measurementName,
	)
	res, err := is.queryDB(stmt)
	if err != nil {
		return time.Time{}, nil
	}
	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values[0]) == 0 {
		is.sugar.Infow("no result returned for last record query", "res", res)
		return time.Time{}, nil
	}
	idxs, err := schema.NewFieldsRegistrar(res[0].Series[0].Columns)
	if err != nil {
		return time.Time{}, nil
	}
	return influxdb.GetTimeFromInterface(res[0].Series[0].Values[0][idxs[schema.Time]])
}
