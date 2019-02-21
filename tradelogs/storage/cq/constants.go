package cq

const (
	hourResampleInterval = "1h"
	hourResampleFor      = "3h"
	dayResampleInterval  = "1h"
	dayResampleFor       = "2d"
	ethwethExcludingTemp = `({{.SrcAddr}}!='{{.ETHTokenAddr}}' OR {{.DstAddr}}!='{{.WETHTokenAddr}}') AND ({{.SrcAddr}}!='{{.WETHTokenAddr}}' OR {{.DstAddr}}!='{{.ETHTokenAddr}}')`
)
