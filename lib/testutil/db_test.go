package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustNewDevelopmentDB(t *testing.T) {
	db, fn := MustNewDevelopmentDB()
	require.NotNil(t, db)
	require.NotNil(t, fn)
	require.NoError(t, fn())
}
