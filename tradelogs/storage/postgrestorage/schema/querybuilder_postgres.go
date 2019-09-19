package schema

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
)

var DateFunctionParams = map[string]string{
	"h": "hour",
	"d": "day",
}

const (
	ethWETHExcludingTmpl = `( NOT EXISTS (SELECT NULL FROM token WHERE address = '{{.ETHTokenAddr}}' AND src_address_id != id )` +
		` OR NOT EXISTS (SELECT NULL FROM token WHERE address = '{{.WETHTokenAddr}}' AND dst_address_id != id ))` +
		` AND ( NOT EXISTS (SELECT NULL FROM token WHERE address = '{{.WETHTokenAddr}}' AND src_address_id != id )` +
		` OR NOT EXISTS (SELECT NULL FROM token WHERE address = '{{.ETHTokenAddr}}' AND dst_address_id != id ))`
)

func BuildDateTruncField(dateTruncParam string, timeZone int8) string {
	if timeZone != 0 && dateTruncParam == "day" {
		var intervalParse = fmt.Sprintf("interval '%d hour'", timeZone)
		return "date_trunc('" + dateTruncParam + "', timestamp + " + intervalParse + ") - " + intervalParse
	}
	return `date_trunc('` + dateTruncParam + `', timestamp)`
}

// RoundTime returns time is rounded by day or hour
// if time is rounded by day, it also use time zone.
func RoundTime(t time.Time, freq string, timeZone int8) time.Time {
	if freq == "hour" {
		return t.Truncate(time.Hour)
	}
	// for example: 6h UTC we want to truncate to 0h in GMT +7
	// 6h UTC  -> 13h GMT +7 -> 0h GMT +7 -> 17h UTC (the day before)
	return t.Add(time.Duration(timeZone) * time.Hour).Truncate(time.Hour * 24).Add(time.Duration(-timeZone) * time.Hour)
}

// BuildEthWethExcludingCondition creates a condition that filter eth-weth trades
func BuildEthWethExcludingCondition() (string, error) {
	var resultBuffer bytes.Buffer

	tpl, err := template.New("exclude eth template").Parse(ethWETHExcludingTmpl)
	if err != nil {
		return "", nil
	}
	err = tpl.Execute(&resultBuffer, struct {
		ETHTokenAddr  string
		WETHTokenAddr string
	}{
		ETHTokenAddr:  blockchain.ETHAddr.Hex(),
		WETHTokenAddr: blockchain.WETHAddr.Hex(),
	})

	if err != nil {
		return "", nil
	}
	return resultBuffer.String(), nil
}
