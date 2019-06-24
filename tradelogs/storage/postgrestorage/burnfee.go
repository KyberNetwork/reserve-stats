package postgrestorage

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

/*// Get aggregated Burn fee by hour or day
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

	if dateFunctionParam, ok = schema.DateFunctionParams[strings.ToLower(freq)]; !ok {
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
}*/

func (tldb *TradeLogDB) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	var (
		timeField string
		err       error
	)
	logger := tldb.sugar.With("from", from, "to", to, "func",
		"tradelogs/storage/postgresql/TradeLogDB.GetAggregatedBurnFee")

	switch strings.ToLower(freq) {
	case "h":
		timeField = schema.BuildDateTruncField("hour", 0)
		from = schema.RoundTime(from, "hour", 0)
		to = schema.RoundTime(to, "hour", 0).Add(time.Hour)
	case "d":
		timeField = schema.BuildDateTruncField("day", 0)
		from = schema.RoundTime(from, "day", 0)
		to = schema.RoundTime(to, "day", 0).Add(time.Hour * 24)

	default:
		return nil, fmt.Errorf("frequency not supported: %v", freq)
	}

	hexAddrs := make([]string, 0)
	for _, rsvAddr := range reserveAddrs {
		hexAddrs = append(hexAddrs, rsvAddr.Hex())
	}

	addrCondition := ""
	if len(hexAddrs) != 0 {
		addrCondition = " AND c.address = ANY($3)"
	}

	integrationQuery := `
		SELECT time, address , SUM(amount) as amount
		FROM (
			SELECT ` + timeField + ` as time, src_burn_amount AS amount, c.address AS address
			FROM "` + schema.TradeLogsTableName + `" b
			INNER JOIN "` + schema.ReserveTableName + `" c ON b.src_reserve_address_id=c.id
			WHERE timestamp >= $1 AND timestamp < $2 ` + addrCondition + `
		UNION ALL
			SELECT ` + timeField + ` as time, dst_burn_amount AS amount, c.address AS address
			FROM "` + schema.TradeLogsTableName + `" b
			INNER JOIN "` + schema.ReserveTableName + `" c ON b.dst_reserve_address_id=c.id
			WHERE timestamp >= $1 AND timestamp < $2 ` + addrCondition + `
		) a GROUP BY time,address
	`

	var records []struct {
		Amount  float64   `db:"amount"`
		Address string    `db:"address"`
		Time    time.Time `db:"time"`
	}

	logger.Debugw("prepare statement", "stmt", integrationQuery)
	err = tldb.db.Select(&records, integrationQuery, from.UTC().Format(schema.DefaultDateFormat),
		to.UTC().Format(schema.DefaultDateFormat), pq.Array(hexAddrs))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}

	result := make(map[ethereum.Address]map[string]float64)
	for _, burnFee := range records {
		reserve := ethereum.HexToAddress(burnFee.Address)
		key := strconv.FormatUint(timeutil.TimeToTimestampMs(burnFee.Time), 10)
		_, ok := result[reserve]
		if !ok {
			result[reserve] = make(map[string]float64)
		}
		result[reserve][key] += burnFee.Amount
	}
	return result, nil

}
