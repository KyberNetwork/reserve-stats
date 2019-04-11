package storage

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestStorage(t *testing.T) {
	var (
		testApplicationName = "test_app_name"
		testAddress1        = "0x64F70539776f08C5EF505254C2426F3e47A5204A"
		testAddress2        = "0xB868636A18c9935D9B259228851cC49245ae68A2"
	)

	sugar := testutil.MustNewDevelopmentSugaredLogger()
	db, fn := testutil.MustNewRandomDevelopmentDB()
	defer func() { assert.NoError(t, fn()) }()

	s, err := NewAppNameDB(sugar, db)
	require.NoError(t, err)

	// refuse to create application with empty application name
	_, _, err = s.CreateOrUpdate(common.Application{})
	require.Error(t, err)

	id, updated, err := s.CreateOrUpdate(common.Application{
		Name:      testApplicationName,
		Addresses: nil,
	})
	require.NoError(t, err)
	t.Logf("application with id %d created", id)
	assert.False(t, updated)

	app, err := s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, app.ID)
	assert.Equal(t, testApplicationName, app.Name)
	assert.Nil(t, app.Addresses)

	id, updated, err = s.CreateOrUpdate(common.Application{
		Name: testApplicationName,
		Addresses: []ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		},
	})
	require.NoError(t, err)
	assert.True(t, updated)

	app, err = s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, app.ID)
	assert.Equal(t, testApplicationName, app.Name)
	assert.ElementsMatch(t,
		[]ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		}, app.Addresses)

	id, updated, err = s.CreateOrUpdate(common.Application{
		ID:   id,
		Name: testApplicationName,
		Addresses: []ethereum.Address{
			ethereum.HexToAddress(testAddress1),
		},
	})
	require.NoError(t, err)
	assert.True(t, updated)

	app, err = s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, app.ID)
	assert.Equal(t, testApplicationName, app.Name)
	assert.ElementsMatch(t,
		[]ethereum.Address{
			ethereum.HexToAddress(testAddress1),
		}, app.Addresses)

	err = s.Update(common.Application{
		ID: id,
		Addresses: []ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		},
	})
	require.NoError(t, err)

	app, err = s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, app.ID)
	assert.Equal(t, testApplicationName, app.Name)
	assert.ElementsMatch(t,
		[]ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		}, app.Addresses)

	err = s.Update(common.Application{
		ID:   id,
		Name: "updated_test_app_name",
	})
	require.NoError(t, err)

	app, err = s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, id, app.ID)
	assert.Equal(t, "updated_test_app_name", app.Name)
	assert.ElementsMatch(t,
		[]ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		}, app.Addresses)

	err = s.Update(common.Application{
		ID:   id + 1000, // not exists
		Name: "updated_test_app_name_2",
	})
	require.Equal(t, ErrNotExists, err)

	apps, err := s.GetAll()
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Application{
		{
			ID:   id,
			Name: "updated_test_app_name",
			Addresses: []ethereum.Address{
				ethereum.HexToAddress(testAddress1),
				ethereum.HexToAddress(testAddress2),
			},
		},
	}, apps)

	apps, err = s.GetAll(WithActiveFilter())
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Application{
		{
			ID:   id,
			Name: "updated_test_app_name",
			Addresses: []ethereum.Address{
				ethereum.HexToAddress(testAddress1),
				ethereum.HexToAddress(testAddress2),
			},
		},
	}, apps)

	apps, err = s.GetAll(WithInactiveFilter())
	require.NoError(t, err)
	assert.Len(t, apps, 0)

	apps, err = s.GetAll(WithNameFilter("updated_test_app_name"))
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Application{
		{
			ID:   id,
			Name: "updated_test_app_name",
			Addresses: []ethereum.Address{
				ethereum.HexToAddress(testAddress1),
				ethereum.HexToAddress(testAddress2),
			},
		},
	}, apps)

	apps, err = s.GetAll(WithNameFilter("random_not_exists_name"))
	require.NoError(t, err)
	assert.Len(t, apps, 0)

	apps, err = s.GetAll(WithNameFilter("random_not_exists_name"), WithActiveFilter())
	require.NoError(t, err)
	assert.Len(t, apps, 0)

	app, err = s.Get(id)
	require.NoError(t, err)
	assert.Equal(t, common.Application{
		ID:   id,
		Name: "updated_test_app_name",
		Addresses: []ethereum.Address{
			ethereum.HexToAddress(testAddress1),
			ethereum.HexToAddress(testAddress2),
		},
	}, app)

	app, err = s.Get(id + 1000)
	require.Equal(t, ErrNotExists, err)

	err = s.Delete(id + 1000)
	require.Equal(t, ErrNotExists, err)

	err = s.Delete(id)
	require.NoError(t, err)
	app, err = s.Get(id)
	require.Equal(t, ErrNotExists, err)
	apps, err = s.GetAll(WithInactiveFilter())
	require.NoError(t, err)
	assert.ElementsMatch(t, []common.Application{
		{
			ID:   id,
			Name: "updated_test_app_name",
			Addresses: []ethereum.Address{
				ethereum.HexToAddress(testAddress1),
				ethereum.HexToAddress(testAddress2),
			},
		},
	}, apps)
}
