package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
)

// Interface is the common interface of users persistent storage.
type Interface interface {
	CreateOrUpdate(common.UserData) error
	// Is KYCed at time return if the user is already KYCed before that time point.
	IsKYCedAtTime(string, time.Time) (bool, error)
}
