package httputil

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	defaultMaxTimeFrame = time.Hour * 24
	defaultTimeFrame    = time.Hour
	defaultFreq         = "d"
)

var (
	defaultFreqs = map[string]time.Duration{
		"h": time.Hour * 24 * 180,     // 180 days
		"d": time.Hour * 24 * 365 * 3, // ~ 3 years
	}
)

// TimeRangeQuery is the common time range query parameters of reporting HTTP APIs.
// The time format is unix milliseconds.
type TimeRangeQuery struct {
	From uint64 `form:"from"`
	To   uint64 `form:"to"`

	// maxTimeFrame is maximum duration allowed between from and to query.
	maxTimeFrame     time.Duration
	defaultTimeFrame time.Duration
}

// TimeRangeQueryValidationOption is the option to configure validation behaviour.
type TimeRangeQueryValidationOption func(q *TimeRangeQuery)

// TimeRangeQueryWithMaxTimeFrame sets the given duration as the max time
// frame when doing validation.
func TimeRangeQueryWithMaxTimeFrame(duration time.Duration) TimeRangeQueryValidationOption {
	return func(q *TimeRangeQuery) {
		q.maxTimeFrame = duration
	}
}

// Validate validates the given time range query and converts them to native
// time.Time format if there is no error.
func (q *TimeRangeQuery) Validate(options ...TimeRangeQueryValidationOption) (time.Time, time.Time, error) {
	for _, option := range options {
		option(q)
	}

	if q.maxTimeFrame == 0 {
		q.maxTimeFrame = defaultMaxTimeFrame
	}

	if q.defaultTimeFrame == 0 {
		q.defaultTimeFrame = defaultTimeFrame
	}

	if q.To == 0 {
		now := time.Now().UTC()
		q.To = timeutil.TimeToTimestampMs(now)
		if q.From == 0 {
			q.From = timeutil.TimeToTimestampMs(now.Add(-q.defaultTimeFrame))
		}
	}

	if q.To < q.From {
		return time.Time{}, time.Time{}, fmt.Errorf("to parameter %d must bigger than from %d parameter", q.To, q.From)
	}

	from := timeutil.TimestampMsToTime(q.From)
	to := timeutil.TimestampMsToTime(q.To)

	if to.Sub(from) > q.maxTimeFrame {
		return time.Time{}, time.Time{}, fmt.Errorf("max time frame exceed, allowed: %s, query: %s",
			q.maxTimeFrame.String(),
			to.Sub(from),
		)
	}

	return from, to, nil
}

// TimeRangeQueryFreq is the same as TimeRangeQuery with additional freq parameter.
type TimeRangeQueryFreq struct {
	TimeRangeQuery
	Freq string `form:"freq"`

	freqs map[string]time.Duration
}

// TimeRangeQueryFreqValidationOption is the option to configure validation behaviour.
type TimeRangeQueryFreqValidationOption func(q *TimeRangeQueryFreq)

// TimeRangeQueryFreqWithValidFreqs configures the validation to use custom valid
// frequencies instead of the default one.
func TimeRangeQueryFreqWithValidFreqs(freqs map[string]time.Duration) TimeRangeQueryFreqValidationOption {
	return func(q *TimeRangeQueryFreq) {
		q.freqs = freqs
	}
}

// Validate validates the given time range query based on given freq.
func (q *TimeRangeQueryFreq) Validate(options ...TimeRangeQueryFreqValidationOption) (time.Time, time.Time, error) {
	for _, option := range options {
		option(q)
	}

	if q.freqs == nil {
		q.freqs = defaultFreqs
	}

	if q.Freq == "" {
		q.Freq = defaultFreq
	}

	maxTimeFrame, ok := q.freqs[q.Freq]
	if !ok {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid frequency: %s", q.Freq)
	}

	from, to, err := q.TimeRangeQuery.Validate(TimeRangeQueryWithMaxTimeFrame(maxTimeFrame))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return from, to, nil
}
