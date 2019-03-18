package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const maxTimeGap = 30000 // 30 secs

//ErrNonceNotInRange error when nonce not in aceptable range
var ErrNonceNotInRange = errors.New("nonce submit is not in aceptable range")

// NonceValidator checking validate by time range
type NonceValidator struct {
	// TimeGap is max time different between client submit timestamp
	// and server time that considered valid. The time precision is millisecond.
	TimeGap uint64
}

// NewNonceValidator return NonceValidator with default value (30 second)
func NewNonceValidator() *NonceValidator {
	return &NonceValidator{
		TimeGap: maxTimeGap,
	}
}

// Validate return error when checking if header date is valid or not
func (v *NonceValidator) Validate(r *http.Request) error {
	nonce, err := strconv.ParseUint(r.Header.Get("nonce"), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse nonce in header. Error: %s", err.Error())
	}
	serverTime := timeutil.UnixMilliSecond()
	start := serverTime - v.TimeGap
	stop := serverTime + v.TimeGap
	if nonce < start || nonce > stop {
		return ErrNonceNotInRange
	}
	return nil
}
