package storage

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	burnVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/burnfee_volume"
	ethereum "github.com/ethereum/go-ethereum/common"
)

var freqToMeasurement = map[string]string{
	"h": "burn_fee_hour",
	"d": "burn_fee_day",
}

// GetAggregatedBurnFee get aggregated burn fee in a time range given the reserve address
func (is *InfluxStorage) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	var (
		measurement string
		addrsStrs   []string
	)
	logger := is.sugar.With("from", from, "to", to, "freq", freq, "reserveAddrs", reserveAddrs)

	measurement, ok := freqToMeasurement[strings.ToLower(freq)]
	if !ok {
		return nil, fmt.Errorf("invalid burn fee frequency %s", freq)
	}
	// The {{range}} in tmpl can't accept other Parameter inside its struct,
	// hence the reserveAddr field name has to be inserted at the tmpl, not to be rendered
	var queryTmpl = `SELECT {{.Fields}} FROM "{{.Measurement}}" WHERE '{{.From }}' <= time AND time <= '{{.To}}' ` +
		`{{if len .Addrs}}AND ({{range $index, $element := .Addrs}}` +
		burnVolumeSchema.ReserveAddr.String() + ` = '{{$element}}'{{if ne $index $.AddrsLastIndex}} OR {{end}}{{end}}){{end}}`

	logger.Debugw("before rendering query statement from template", "query_tempalte", queryTmpl)
	tmpl, err := template.New("queryStmt").Parse(queryTmpl)
	if err != nil {
		return nil, err
	}

	var queryStmtBuf bytes.Buffer
	for _, rsvAddr := range reserveAddrs {
		addrsStrs = append(addrsStrs, rsvAddr.Hex())
	}
	logger.Debug(burnVolumeSchema.ReserveAddr.String())
	if err = tmpl.Execute(&queryStmtBuf, struct {
		Fields         string
		Measurement    string
		From           string
		To             string
		Addrs          []string
		AddrsLastIndex int
	}{
		Fields:         burnVolumeSchema.SumAmount.String() + `, ` + burnVolumeSchema.ReserveAddr.String(),
		Measurement:    measurement,
		From:           from.Format(time.RFC3339),
		To:             to.Format(time.RFC3339),
		Addrs:          addrsStrs,
		AddrsLastIndex: len(reserveAddrs) - 1,
	}); err != nil {
		return nil, err
	}

	logger.Debugw("rendered query statement", "rendered_template", queryStmtBuf.String())

	res, err := is.queryDB(is.influxClient, queryStmtBuf.String())
	if err != nil {
		return nil, err
	}

	if len(res[0].Series) == 0 {
		logger.Debug("empty aggregated burn fee result")
		return nil, nil
	}

	result := make(map[ethereum.Address]map[string]float64)

	idxs, err := burnVolumeSchema.NewFieldsRegistrar(res[0].Series[0].Columns)
	if err != nil {
		return nil, err
	}
	for _, value := range res[0].Series[0].Values {
		ts, amount, reserve, err := is.rowToAggregatedBurnFee(value, idxs)
		if err != nil {
			return nil, err
		}

		key := strconv.FormatUint(timeutil.TimeToTimestampMs(ts), 10)

		_, ok := result[reserve]
		if !ok {
			result[reserve] = make(map[string]float64)
		}
		result[reserve][key] = amount
	}

	return result, nil
}
