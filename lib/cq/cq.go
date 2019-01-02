package cq

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	cqTemplate = `{{if .IsCQ}}CREATE CONTINUOUS QUERY "{{.ActualName}}" on "{{.Database}}" ` +
		`{{if or .ResampleEveryInterval .ResampleForInterval}}RESAMPLE {{if .ResampleEveryInterval}}EVERY {{.ResampleEveryInterval}} {{end}}{{if .ResampleForInterval}}FOR {{.ResampleForInterval}} {{end}}{{end}}` +
		`BEGIN {{end}}{{.ActualQuery}}` +
		`{{if not .GroupByQuery}} GROUP BY {{else}}, {{end}}time({{.TimeInterval}}{{if .OffsetInterval}},{{.OffsetInterval}}{{end}}) {{if .IsCQ}}END{{end}}`
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
	queries, err := cq.prepareQueries(true)
	if err != nil {
		return nil, err
	}
	cq.queries = queries
	queries, err = cq.prepareQueries(false)
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

// hasGroupBy decide if the query has a GROUP BY from the main select
// the logic is to count the number of complete bracket from the last GROUP BY
func hasGroupBy(query string) bool {
	idx := strings.LastIndex(query, "GROUP BY")
	if idx < 0 {
		return false
	}
	totalBracket := 0
	for i := idx; i < len(query); i++ {
		if query[i] == ')' {
			totalBracket--
		}
		if query[i] == '(' {
			totalBracket++
		}
	}
	if totalBracket == 0 {
		return true
	}
	return false
}

func (cq *ContinuousQuery) prepareQueries(isCQ bool) ([]string, error) {
	var queries []string

	cq.OffsetIntervals = append(cq.OffsetIntervals, "")
	for _, offsetInterval := range cq.OffsetIntervals {
		var query bytes.Buffer

		tmpl, err := template.New("cq.prepareQueries").Parse(cqTemplate)
		if err != nil {
			return nil, err
		}
		// Modify name and query so that mutiple queries are allowed
		actualName := cq.Name
		actualQuery := cq.Query
		if offsetInterval != "" {
			// assume that the offsetInterval value is set correctly
			if strings.Contains(offsetInterval, "-") {
				offsetIntervalName := strings.Replace(offsetInterval, "-", "", 1)
				actualName = actualName + "_" + offsetIntervalName
				actualQuery = modifyINTOclause(cq.Query, offsetIntervalName)
			} else {
				actualName = actualName + "_" + "minus" + offsetInterval
				actualQuery = modifyINTOclause(cq.Query, offsetInterval)
			}
		}
		err = tmpl.Execute(&query, struct {
			*ContinuousQuery
			ActualName     string
			ActualQuery    string
			GroupByQuery   bool // whether the query included GROUP BY statement
			IsCQ           bool // whether the query is a continous query or not
			OffsetInterval string
		}{
			ContinuousQuery: cq,
			ActualName:      actualName,
			ActualQuery:     actualQuery,
			GroupByQuery:    hasGroupBy(cq.Query),
			IsCQ:            isCQ,
			OffsetInterval:  offsetInterval,
		})
		if err != nil {
			return nil, err
		}

		queries = append(queries, query.String())
	}

	return queries, nil
}

// GetCurrentCQs return a map[cqName]query of current CQs in the database
func (cq *ContinuousQuery) GetCurrentCQs(c client.Client, sugar *zap.SugaredLogger) (map[string]string, error) {
	var (
		qr     = "SHOW CONTINUOUS QUERIES"
		result = make(map[string]string)
	)
	resp, err := queryDB(c, qr, cq.Database)
	if err != nil {
		return result, err
	}
	if len(resp[0].Series) == 0 {
		sugar.Debugw("current cqs request: empty result")
		return result, nil
	}
	for _, row := range resp[0].Series {
		if row.Name == cq.Database {
			for _, v := range row.Values {
				cqName, ok := v[0].(string)
				if !ok {
					return result, fmt.Errorf("can not decode cqName value %v to string", v[0])
				}
				cq, ok := v[1].(string)
				if !ok {
					return result, fmt.Errorf("can not decode cq value %v to string", v[1])
				}
				result[cqName] = cq
			}
		}
	}
	return result, nil
}

// Deploy ensures that all configured cqs are deployed in given InfluxDB server. This method is safe to run multiple
// times without changes as it will checks if the CQ exists and updated first.
// client and logger will be need to inject from higher level call
func (cq *ContinuousQuery) Deploy(c client.Client, sugar *zap.SugaredLogger) error {
	for _, query := range cq.queries {
		sugar.Debugw("Executing Query", "query", query)

		_, err := queryDB(c, query, cq.Database)
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute runs all CQs query to aggregate historical data.
// This method is intended to use once in new deployment.
func (cq *ContinuousQuery) Execute(c client.Client, sugar *zap.SugaredLogger) error {
	for _, query := range cq.historicalQueries {
		sugar.Debugw("Executing Query", "query", query)
		_, err := queryDB(c, query, cq.Database)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drop drop the cq from its database. Since influxDB doesn't raise error if cq doesn't exist,
// This won't check if the cq is already there.
func (cq *ContinuousQuery) Drop(c client.Client, sugar *zap.SugaredLogger) error {
	if len(cq.OffsetIntervals) == 0 {
		cq.OffsetIntervals = []string{""}
	}
	for _, offset := range cq.OffsetIntervals {
		name := cq.Name
		if offset != "" {
			name = cq.Name + "_" + offset
		}
		sugar.Debugw("Drop cq", "cq name", name)
		if _, err := cq.queryDB(c, fmt.Sprintf("DROP CONTINUOUS QUERY %s ON %s", name, cq.Database)); err != nil {
			return err
		}
	}
	return nil
}

// queryDB convenience function to query the database
func (cq *ContinuousQuery) queryDB(c client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: cq.Database,
	}
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
