package postgrestorage

import (
	"bytes"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"html/template"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var dateFunctionParams = map[string]string{
	"h": "hour",
	"d": "day",
}

// Get aggregated Burn fee by hour or day
func (tldb *TradeLogDB) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	var (
		dateFunctionParam string
		addressHexs       []string
		ok                bool
	)

	for _, rsvAddr := range reserveAddrs {
		addressHexs = append(addressHexs, rsvAddr.Hex())
	}
	logger := tldb.sugar.With("from", from, "to", to, "freq", freq, "reserveAddrs", reserveAddrs)

	if dateFunctionParam, ok = dateFunctionParams[strings.ToLower(freq)]; !ok {
		return nil, fmt.Errorf("invalid burn fee frequency %s", freq)
	}

	var queryTmp = `SELECT SUM({{.AmountColumn}}) AS amount, b.{{.AddressColumn}} AS address,` +
		` date_trunc('{{.DateParam}}',{{.TimeColumn}}) as time ` +
		` FROM {{.TradeLogsTableName}} AS a{{if .HasAddr}}, {{.ReserveTableName}} AS b{{end}}` +
		` WHERE date_trunc('{{.DateParam}}',{{.TimeColumn}}) >= $1 AND $2 >= date_trunc('{{.DateParam}}',{{.TimeColumn}})` +
		`{{if .HasAddr}} AND b.{{.AddressColumn}} = ANY($3){{end}}` +
		` AND  a.{{.AddressIdColumn}}=b.{{.MatchIdColumn}}` +
		` GROUP BY address, time`

	logger.Debugw("before rendering query statement from template", "query_template", queryTmp)
	tpl, err := template.New("burn fee query template").Parse(queryTmp)
	if err != nil {
		return nil, err
	}

	result := make(map[ethereum.Address]map[string]float64)

	for _, param := range []struct {
		AmountColumn    string
		AddressIdColumn string
	}{{AmountColumn: "src_burn_amount", AddressIdColumn: "src_reserve_address_id"},
		{AmountColumn: "dst_burn_amount", AddressIdColumn: "dst_reserve_address_id"}} {
		var queryStmtBuf bytes.Buffer
		if err = tpl.Execute(&queryStmtBuf, struct {
			DateParam          string
			AmountColumn       string
			AddressColumn      string
			TimeColumn         string
			AddressIdColumn    string
			MatchIdColumn      string
			TradeLogsTableName string
			ReserveTableName   string
			HasAddr            bool
		}{
			DateParam:          dateFunctionParam,
			AmountColumn:       param.AmountColumn,
			AddressColumn:      "address",
			TimeColumn:         "timestamp",
			AddressIdColumn:    param.AddressIdColumn,
			MatchIdColumn:      "id",
			TradeLogsTableName: schema.TradeLogsTableName,
			ReserveTableName:   schema.ReserveTableName,
			HasAddr:            len(reserveAddrs) != 0,
		}); err != nil {
			return nil, err
		}

		logger.Debugw("prepare statement", "smt", queryStmtBuf.String())

		var burnFees []struct {
			Amount  float64   `db:"amount"`
			Address string    `db:"address"`
			Time    time.Time `db:"time"`
		}

		if err = tldb.db.Select(&burnFees, queryStmtBuf.String(), from.UTC().Format(schema.DefaultDateFormat),
			to.UTC().Format(schema.DefaultDateFormat), pq.Array(addressHexs)); err != nil {
			return nil, err
		}

		for _, burnFee := range burnFees {
			reserve := ethereum.HexToAddress(burnFee.Address)
			key := strconv.FormatUint(timeutil.TimeToTimestampMs(burnFee.Time), 10)
			_, ok := result[reserve]
			if !ok {
				result[reserve] = make(map[string]float64)
			}
			result[reserve][key] += burnFee.Amount
		}

	}

	return result, nil
}
