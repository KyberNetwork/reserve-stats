package http

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func newTestServer(rts storage.ReserveTransactionStorage, sugar *zap.SugaredLogger) (*Server, error) {
	return NewServer(
		sugar,
		"",
		rts,
	), nil

}

func TestERC20Transfer(t *testing.T) {
	_, storage := testutil.MustNewDevelopmentDB()
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	rts, err := postgres.NewStorage(sugar, storage)
	require.NoError(t, err)

	s, err := newTestServer(rts, sugar)
	require.NoError(t, err)
	s.register()

	defer func(t *testing.T) {
		require.NoError(t, rts.TearDown())
	}(t)

	err = rts.StoreReserve(ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"), common.CompanyWallet.String())
	require.NoError(t, err)

	// prepare data
	testWalletERC20Transfer := []common.ERC20Transfer{
		{
			Timestamp:       timeutil.TimestampMsToTime(1554094535000),
			Hash:            ethereum.HexToHash("0xf18cc8635570d4be2ec39a89234219bce64785333978c56f98c6772a9a200942"),
			From:            ethereum.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"),
			ContractAddress: ethereum.HexToAddress("0x255Aa6DF07540Cb5d3d297f0D0D4D84cb52bc8e6"),
			To:              ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
			BlockNumber:     7480572,
			Value:           big.NewInt(1634062936237169),
			Gas:             376367,
			GasUsed:         321367,
			GasPrice:        big.NewInt(4179327936),
		},
		{
			Timestamp:       timeutil.TimestampMsToTime(1554199276000),
			Hash:            ethereum.HexToHash("0x5a96c395432dcc2d416d4dbe069449381edae39d24fc9e7dda15148ac7b3fdd7"),
			From:            ethereum.HexToAddress("0x44d34A119BA21A42167FF8B77a88F0Fc7BB2Db90"),
			ContractAddress: ethereum.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"),
			To:              ethereum.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"),
			BlockNumber:     7480583,
			Value:           big.NewInt(11321854),
			Gas:             338021,
			GasUsed:         191589,
			GasPrice:        big.NewInt(6000000000),
		},
	}

	err = rts.StoreERC20Transfer(testWalletERC20Transfer, ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"))
	require.NoError(t, err)

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test without wallet address and token address",
			Endpoint: "/wallet/transactions?from=1554094535000&to=1554199276001",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.ERC20Transfer
				sugar.Debugw("resutl", "response", resp.Body)
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
				assert.Equal(t, 2, len(result))
			},
		},
		{
			Msg:      "Test wallet address",
			Endpoint: "/wallet/transactions?from=1554094535000&to=1554199276001&wallet=0x63825c174ab367968EC60f061753D3bbD36A0D8F",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.ERC20Transfer
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
				assert.Equal(t, 1, len(result))
			},
		},
		{
			Msg:      "Test token address",
			Endpoint: "/wallet/transactions?from=1554094535000&to=1554199276001&token=0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.ERC20Transfer
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
				assert.Equal(t, 1, len(result))
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
