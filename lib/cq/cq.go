package cq

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	cqTemplate = `CREATE CONTINUOUS QUERY "{{.Name}}" on "{{.Database}}" ` +
		`{{if or .ResampleEveryInterval .ResampleForInterval}}RESAMPLE {{if .ResampleEveryInterval}}EVERY {{.ResampleEveryInterval}} {{end}}{{if .ResampleForInterval}}FOR {{.ResampleForInterval}} {{end}}{{end}}` +
		`BEGIN {{.Query}}` +
		`{{if not .GroupByQuery}} GROUP BY {{else}}, {{end}}time({{.TimeInterval}}{{if .OffsetInterval}},{{.OffsetInterval}}{{end}}) END`
	hqTemplate = `{{.Query}}` + `{{if not .GroupByQuery}} GROUP BY {{else}}, {{end}}time({{.TimeInterval}}{{if .OffsetInterval}},{{.OffsetInterval}}{{end}})`
)

// NewContinuousQuery creates new ContinuousQuery instance.
func NewContinuousQuery(
	name, database, resampleEveryInterval, resampleForInterval, query,
	timeInterval string, offsetIntervals []string) (*ContinuousQuery, error) {
	cq := &ContinuousQuery{
		Name:                  name,
		Database:              database,
		ResampleEveryInterval: resampleEveryInterval,
		ResampleForInterval:   resampleForInterval,
		Query:                 query,
		TimeInterval:          timeInterval,
		OffsetIntervals:       offsetIntervals,
	}
	queries, err := cq.prepareQueries(cqTemplate)
	if err != nil {
		return nil, err
	}
	queries, err = cq.prepareQueries(hqTemplate)
	if err != nil {
		return nil, err
	}
	cq.historicalQueries = queries
	return cq, nil
}

// ContinuousQuery represents an InfluxDB Continuous Query.
// By design ContinuousQuery doesn't try to be smart, it does not attempt to parse/validate any field,
// just act as a templating engine.
//
// Example:
// CREATE CONTINUOUS QUERY <cq_name> ON <database_name>
// RESAMPLE EVERY <interval> FOR <interval>
// BEGIN
// <cq_query>
// END
type ContinuousQuery struct {
	Name                  string
	Database              string
	ResampleEveryInterval string
	ResampleForInterval   string
	// the Query string without the GROUP BY time part which will be added by
	// examining TimeInterval and OffsetIntervals.
	Query           string
	TimeInterval    string
	OffsetIntervals []string

	queries           []string
	historicalQueries []string
}

func modifyINTOclause(query string, suffix string) string {
	words := strings.Fields(query)
	if len(words) == 0 {
		return query
	}
	for i, word := range words {
		if (strings.ToUpper(strings.TrimSpace(word)) == "INTO") && (i+1 < len(words)) {
			words[i+1] = words[i+1] + "_" + suffix
			break
		}
	}
	var s string
	for _, word := range words {
		s = s + word + " "
	}
	return strings.TrimSuffix(s, " ")
}

func (cq *ContinuousQuery) prepareQueries(tmp string) ([]string, error) {
	var queries []string

	if len(cq.OffsetIntervals) == 0 {
		cq.OffsetIntervals = []string{""}
	}

	for _, offsetInterval := range cq.OffsetIntervals {
		var query bytes.Buffer

		tmpl, err := template.New("cq.prepareQueries").Parse(tmp)
		if err != nil {
			return nil, err
		}

		//create a forkedCq to alter the name, otherwise the cqs will failed when their names duplicated.
		forkedCq := ContinuousQuery{}
		forkedCq = *cq
		forkedCq.Name = offsetInterval + forkedCq.Name
		forkedCq.Query = modifyINTOclause(forkedCq.Query, offsetInterval)
		err = tmpl.Execute(&query, struct {
			*ContinuousQuery
			GroupByQuery   bool // whether the query included GROUP BY statement
			OffsetInterval string
		}{
			ContinuousQuery: &forkedCq,
			GroupByQuery:    strings.Contains(cq.Query, "GROUP BY"),
			OffsetInterval:  offsetInterval,
		})
		if err != nil {
			return nil, err
		}

		queries = append(queries, query.String())
	}

	return queries, nil
}

// Deploy ensures that all configured cqs are deployed in given InfluxDB server. This method is safe to run multiple
// times without changes as it will checks if the CQ exists and updated first.
// TODO: should we move client to constructor?
func (cq *ContinuousQuery) Deploy(c client.Client, sugar *zap.SugaredLogger) error {
	// TODO: implement this
	return nil
}

// Execute runs all CQs query to aggregate historical data.
// This method is intended to use once in new deployment.
func (cq *ContinuousQuery) Execute(c client.Client, sugar *zap.SugaredLogger) error {
	for _, query := range cq.historicalQueries {
		sugar.Debugw("Executing Query", "query", query)
		q := client.Query{
			Command:  query,
			Database: cq.Database,
		}
		response, err := c.Query(q)
		if err != nil {
			return err
		}
		if response.Error() != nil {
			return response.Error()
		}
	}
	return nil
}
