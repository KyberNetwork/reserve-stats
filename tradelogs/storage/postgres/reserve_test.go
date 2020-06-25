package postgres

import (
	"context"
	"math/big"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	alchemyRopsten = "https://eth-ropsten.alchemyapi.io/v2/GvXu_IIrL0U10ZpgTXKjOGIA06KcwzEK"
)

func TestAddUpdateReserve(t *testing.T) {
	const (
		dbName                   = "test_reserve_v4"
		addReserveToStorageEvent = "0x4649526e2876a69a4439244e5d8a32a6940a44a92b5390fdde1c22a26cc54004"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		// require.NoError(t, testStorage.tearDown(dbName))
	}()

	c, err := ethclient.Dial(alchemyRopsten)
	require.NoError(t, err)
	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(addReserveToStorageEvent),
		},
	}
	query := ether.FilterQuery{
		FromBlock: big.NewInt(8008389),
		ToBlock:   big.NewInt(8008389),
		Addresses: []ethereum.Address{ethereum.HexToAddress("0xa4ead31a6c8e047e01ce1128e268c101ad391959")},
		Topics:    topics,
	}
	logs, err := c.FilterLogs(context.Background(), query)
	require.NoError(t, err)
	kyberStorageContract, err := contracts.NewKyberStorage(
		ethereum.HexToAddress("0xa4ead31a6c8e047e01ce1128e268c101ad391959"),
		c,
	)
	require.NoError(t, err)

	for _, log := range logs {
		t.Log("hash", log.TxHash.String())
		d, err := kyberStorageContract.ParseAddReserveToStorage(log)
		require.NoError(t, err)
		err = testStorage.AddUpdateReserve(&tradelogs.AddReserveToStorage{
			Reserve:      d.Reserve,
			ReserveID:    d.ReserveId,
			BlockNumber:  log.BlockNumber,
			RebateWallet: d.RebateWallet,
		})
		require.NoError(t, err)
	}
}
