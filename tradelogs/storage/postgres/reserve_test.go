package postgres

import (
	"context"
	"math/big"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	infuraRopsten = "https://ropsten.infura.io/v3/3243b06bf0334cff8e468bf90ce6e08c"
)

func TestAddUpdateReserve(t *testing.T) {
	const (
		dbName                   = "test_reserve_v4"
		addReserveToStorageEvent = "0x50b2ce9e8f1a63ceaed262cc854dbf741b216e6429f7ba38403afbcdddc7f1ea"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()

	c, err := ethclient.Dial(infuraRopsten)
	require.NoError(t, err)
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(addReserveToStorageEvent),
		},
	}
	query := ether.FilterQuery{
		FromBlock: big.NewInt(8192307),
		ToBlock:   big.NewInt(8192307),
		Addresses: []ethereum.Address{ethereum.HexToAddress("0x688bf5eec43e0799c5b9c1612f625f7b93fe5434")},
		Topics:    topics,
	}
	logs, err := c.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	kyberStorageContract, err := contracts.NewKyberStorage(
		ethereum.HexToAddress("0x688bf5eec43e0799c5b9c1612f625f7b93fe5434"),
		c,
	)
	require.NoError(t, err)

	for _, log := range logs {
		t.Log("hash", log.TxHash.String())
		d, err := kyberStorageContract.ParseAddReserveToStorage(log)
		require.NoError(t, err)
		err = testStorage.saveReserve([]common.Reserve{
			{
				Address:      d.Reserve,
				ReserveID:    d.ReserveId,
				ReserveType:  uint64(d.ReserveType),
				BlockNumber:  log.BlockNumber,
				RebateWallet: d.RebateWallet,
			},
		})
		require.NoError(t, err)
	}
}
