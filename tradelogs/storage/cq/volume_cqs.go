package cq

import (
	"bytes"
	"text/template"

	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
)

const (
	rsvVolTemplate    = `SELECT SUM({{.AmountType}}) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO {{.MeasurementName}} FROM (SELECT {{.AmountType}}, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE ((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) AND {{.RsvAddressType}}!='') GROUP BY {{.AddressType}},{{.RsvAddressType}}`
	rsvVolHourMsmName = `rsv_volume_hour`
	rsvVolDayMsmName  = `rsv_volume_day`
)

// CreateAssetVolumeCqs return a set of cqs required for asset volume aggregation
func CreateAssetVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result []*libcq.ContinuousQuery
	)
	assetVolDstHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE ((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY dst_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstHourCqs)
	assetVolSrcHourCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE ((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY src_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcHourCqs)
	assetVolDstDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_dst_day",
		dbName,
		"1h",
		"2d",
		"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_day FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE ((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY dst_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolDstDayCqs)

	assetVolSrcDayCqs, err := libcq.NewContinuousQuery(
		"asset_volume_src_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_day FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE ((src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))) GROUP BY src_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, assetVolSrcDayCqs)
	return result, nil
}

// RsvFieldsType declare the set of names requires to completed a reserveVolume Cqs
type RsvFieldsType struct {
	AmountType     string
	RsvAddressType string
	AddressType    string
}

func renderRsvCqFromTemplate(tmpl *template.Template, mName string, types RsvFieldsType) (string, error) {
	var query bytes.Buffer
	err := tmpl.Execute(&query, struct {
		RsvFieldsType
		MeasurementName string
	}{
		RsvFieldsType:   types,
		MeasurementName: mName,
	})
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

// CreateReserveVolumeCqs return a set of cqs required for asset volume aggregation
func CreateReserveVolumeCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var (
		result     []*libcq.ContinuousQuery
		cqsGroupBY = map[string]RsvFieldsType{
			"rsv_volume_src_src": RsvFieldsType{
				AmountType:     "src_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_src_dst": RsvFieldsType{
				AmountType:     "src_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "src_addr"},
			"rsv_volume_dst_src": RsvFieldsType{
				AmountType:     "dst_amount",
				RsvAddressType: "src_rsv_addr",
				AddressType:    "dst_addr"},
			"rsv_volume_dst_dst": RsvFieldsType{
				AmountType:     "dst_amount",
				RsvAddressType: "dst_rsv_addr",
				AddressType:    "dst_addr"},
		}
	)

	tpml, err := template.New("cq.CreateReserveVolumeCqs").Parse(rsvVolTemplate)
	if err != nil {
		return nil, err
	}

	for name, types := range cqsGroupBY {
		query, err := renderRsvCqFromTemplate(tpml, rsvVolHourMsmName, types)
		if err != nil {
			return nil, err
		}
		hourCQ, err := libcq.NewContinuousQuery(
			name+"_hour",
			dbName,
			hourResampleInterval,
			hourResampleFor,
			query,
			"1h",
			[]string{},
		)
		if err != nil {
			return nil, err
		}
		result = append(result, hourCQ)

		query, err = renderRsvCqFromTemplate(tpml, rsvVolDayMsmName, types)
		if err != nil {
			return nil, err
		}
		dayCQ, err := libcq.NewContinuousQuery(
			name+"_day",
			dbName,
			dayResampleInterval,
			dayResampleFor,
			query,
			"1d",
			[]string{},
		)
		if err != nil {
			return nil, err
		}
		result = append(result, dayCQ)
	}

	return result, nil
}
