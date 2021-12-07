package common

// OrderState is type of Huobi order states.
//go:generate stringer -type=OrderState -linecomment
type OrderState int

const (
	// PreSubmitted is prepare to submit state
	PreSubmitted OrderState = iota + 1 // pre-submitted
	// Submitting is submitting state
	Submitting // submitting
	// Submitted is submitted state
	Submitted // submitted
	// PartialFilled is partial filled state
	PartialFilled // partial-filled
	// PartialCanceled is partial canceled state
	PartialCanceled // partial-canceled
	// Filled is filled state
	Filled // filled
	// Canceled is canceled state
	Canceled // canceled
	// Failed is failed state
	Failed // failed
)

// OrderStates map a string to OrderState
var OrderStates = map[string]OrderState{
	"pre-submitted":    PreSubmitted,
	"submitting":       Submitting,
	"submitted":        Submitted,
	"partial-filled":   PartialFilled,
	"partial-canceled": PartialCanceled,
	"filled":           Filled,
	"canceled":         Canceled,
	"failed":           Failed,
}
