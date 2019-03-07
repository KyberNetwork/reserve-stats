package storage

import "errors"

var (
	// ErrNotExists is the error returned when querying object does not exist in database.
	ErrNotExists = errors.New("not exists")

	// ErrExists is the error returned when the record is already exists in database.
	ErrExists = errors.New("already exists")
)
