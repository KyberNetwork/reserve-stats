package schema

import (
	"bytes"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"html/template"
	"strings"
	"time"
)

var DateFunctionParams = map[string]string{
	"h": "hour",
	"d": "day",
}

const (
	timeConditionTemplate = `date_trunc('{{.DateParam}}',{{.TimeColumn}}) >= '{{.StartTime}}'` +
		` AND '{{.EndTime}}' >= date_trunc('{{.DateParam}}',{{.TimeColumn}})`
	ethWETHExcludingTmpl = `({{.SrcAddr}}!='{{.ETHTokenAddr}}' OR {{.DstAddr}}!='{{.WETHTokenAddr}}')` +
		` AND ({{.SrcAddr}}!='{{.WETHTokenAddr}}' OR {{.DstAddr}}!='{{.ETHTokenAddr}}')`
)

func BuildTimeCondition(start time.Time, end time.Time, frequency string) (string, error) {
	var (
		dateFunctionParam string
		resultBuffer      bytes.Buffer
		ok                bool
	)
	tpl, err := template.New("time condition template").Parse(timeConditionTemplate)
	if err != nil {
		return "", err
	}

	if dateFunctionParam, ok = DateFunctionParams[strings.ToLower(frequency)]; !ok {
		return "", fmt.Errorf("invalid burn fee frequency %s", frequency)
	}

	err = tpl.Execute(&resultBuffer, struct {
		DateParam  string
		TimeColumn string
		StartTime  string
		EndTime    string
	}{
		DateParam:  dateFunctionParam,
		TimeColumn: "timestamp",
		StartTime:  start.UTC().Format(DefaultDateFormat),
		EndTime:    end.UTC().Format(DefaultDateFormat),
	})
	if err != nil {
		return "", err
	}
	return resultBuffer.String(), nil
}

func BuildEthWethExcludingCondition(srcTable string, dstTable string) (string, error) {
	var resultBuffer bytes.Buffer

	tpl, err := template.New("exclude eth template").Parse(ethWETHExcludingTmpl)
	if err != nil {
		return "", nil
	}
	err = tpl.Execute(&resultBuffer, struct {
		SrcAddr       string
		DstAddr       string
		ETHTokenAddr  string
		WETHTokenAddr string
	}{
		SrcAddr:       srcTable + ".address",
		DstAddr:       dstTable + ".address",
		ETHTokenAddr:  blockchain.ETHAddr.Hex(),
		WETHTokenAddr: blockchain.WETHAddr.Hex(),
	})

	if err != nil {
		return "", nil
	}
	return resultBuffer.String(), nil
}

func BuildAddressCondition(srcTable string, dstTable string, address string) (string, error) {
	return "", nil
}
