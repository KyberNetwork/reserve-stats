package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage/postgresql"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

var (
	tdb *sqlx.DB
	ts  *Server
	tst *postgresql.Storage
	tts time.Time
)

func TestReserveAddressesCreate(t *testing.T) {
	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "create a reserve addresses successfully",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F",
  "type": "reserve",
  "description": "main reserve"
}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusCreated, resp.Code)
				var addr = &common.ReserveAddress{}
				err := json.NewDecoder(resp.Body).Decode(addr)
				require.NoError(t, err)
				assert.NotEmpty(t, addr.ID)

				var stored = &postgresql.ReserveAddress{}
				err = tdb.Get(stored, `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`, addr.ID)
				require.NoError(t, err)
				assert.Equal(t, "0x63825c174ab367968EC60f061753D3bbD36A0D8F", stored.Address)
				assert.Equal(t, common.Reserve.String(), stored.Type)
				assert.Equal(t, "main reserve", stored.Description)
				assert.Equal(t, stored.Timestamp.Time.Unix(), tts.Unix())
			},
		},
		{
			Msg:      "create a reserve address with duplicated Ethereum address",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F",
  "type": "reserve",
  "description": "main reserve"
}`),
			Assert: httputil.AssertCode(http.StatusConflict),
		},
		{
			Msg:      "create a reserve address with missing address",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "type": "reserve",
  "description": "main reserve"
}`),
			Assert: httputil.AssertCode(http.StatusBadRequest),
		},
		{
			Msg:      "create a reserve address with missing type",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "address": "0x63825c174ab367968EC60f061753D3bbD36A0D8F",
  "description": "main reserve"
}`),
			Assert: httputil.AssertCode(http.StatusBadRequest),
		},
		{
			Msg:      "create a reserve address with missing description",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "address": "0x63825c174ab367968EC60f061753D3bbD36A0D80",
  "type": "pricing_operator"
}`),

			Assert: httputil.AssertCode(http.StatusCreated),
		},
		{
			Msg:      "create a reserve address with invalid address",
			Endpoint: "/addresses",
			Method:   http.MethodPost,
			Body: []byte(`{
  "address": "invalid-address",
  "type": "reserve",
  "description": "main reserve"
}`),
			Assert: httputil.AssertCode(http.StatusBadRequest),
		},
	}

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}
}

func TestReserveAddressesGet(t *testing.T) {
	var (
		testAddress     = ethereum.HexToAddress("0x78bf540f3198bc64599ac46b1b43c8012957bdf2e4b2871403a332f7b995da98")
		testDescription = "test pricing operator"
	)

	t.Log("creating a test reserve address")
	id, err := tst.Create(testAddress, common.PricingOperator, testDescription)
	require.NoError(t, err)

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "get a existing reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)
				var addr = &common.ReserveAddress{}
				err := json.NewDecoder(resp.Body).Decode(addr)
				require.NoError(t, err)
				assert.Equal(t, id, addr.ID)
				assert.Equal(t, testAddress, addr.Address)
				assert.Equal(t, common.PricingOperator, addr.Type)
				assert.Equal(t, testDescription, addr.Description)
				assert.Equal(t, tts.UTC().Unix(), addr.Timestamp.Unix())
			},
		},
		{
			Msg:      "get an non existing address",
			Endpoint: fmt.Sprintf("/addresses/%d", id+100),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			Msg:      "get an invalid id",
			Endpoint: fmt.Sprintf("/addresses/%s", "invalid-id"),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}
}

func TestReserveAddressUpdate(t *testing.T) {
	var (
		testAddress     = ethereum.HexToAddress("0x675ADFEcaDe88cE7342BBc34FeF1A1F01CB2a8c4")
		testDescription = "test sanity operator"
	)

	t.Log("creating a test reserve address")
	id, err := tst.Create(testAddress, common.SanityOperator, testDescription)
	require.NoError(t, err)

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "update address of a reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id),
			Method:   http.MethodPut,
			Body: []byte(`{
  "address": "0x818e6fecd516ecc3849daf6845e3ec868087b755"
}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusNoContent, resp.Code)

				var updated = &postgresql.ReserveAddress{}
				err = tdb.Get(updated, `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`, id)
				require.NoError(t, err)
				assert.Equal(t,
					ethereum.HexToAddress("0x818e6fecd516ecc3849daf6845e3ec868087b755"),
					ethereum.HexToAddress(updated.Address))
				assert.Equal(t, common.SanityOperator.String(), updated.Type)
				assert.Equal(t, testDescription, updated.Description)
				assert.Equal(t, updated.Timestamp.Time.Unix(), tts.Unix())
			},
		},
		{
			Msg:      "update description of a reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id),
			Method:   http.MethodPut,
			Body: []byte(`{
  "description": "updated description"
}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusNoContent, resp.Code)

				var updated = &postgresql.ReserveAddress{}
				err = tdb.Get(updated, `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`, id)
				require.NoError(t, err)
				assert.Equal(t,
					ethereum.HexToAddress("0x818e6fecd516ecc3849daf6845e3ec868087b755"),
					ethereum.HexToAddress(updated.Address))
				assert.Equal(t, common.SanityOperator.String(), updated.Type)
				assert.Equal(t, "updated description", updated.Description)
				assert.Equal(t, updated.Timestamp.Time.Unix(), tts.Unix())
			},
		},
		{
			Msg:      "update type of a reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id),
			Method:   http.MethodPut,
			Body: []byte(`{
  "type": "intermediate_operator"
}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusNoContent, resp.Code)

				var updated = &postgresql.ReserveAddress{}
				err = tdb.Get(updated, `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`, id)
				require.NoError(t, err)
				assert.Equal(t,
					ethereum.HexToAddress("0x818e6fecd516ecc3849daf6845e3ec868087b755"),
					ethereum.HexToAddress(updated.Address))
				assert.Equal(t, common.IntermediateOperator.String(), updated.Type)
				assert.Equal(t, "updated description", updated.Description)
				assert.Equal(t, updated.Timestamp.Time.Unix(), tts.Unix())
			},
		},
		{
			Msg:      "update all information of a reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id),
			Method:   http.MethodPut,
			Body: []byte(`{
"address": "0x885Ecc9c31320993e0fE6148cd90F5159DB19d36",
"type": "reserve",
"description": "super reserve 1"
}`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusNoContent, resp.Code)

				var updated = &postgresql.ReserveAddress{}
				err = tdb.Get(updated, `SELECT id, address, type, description, timestamp
FROM addresses
WHERE id = $1`, id)
				require.NoError(t, err)
				assert.Equal(t,
					ethereum.HexToAddress("0x885Ecc9c31320993e0fE6148cd90F5159DB19d36"),
					ethereum.HexToAddress(updated.Address))
				assert.Equal(t, common.Reserve.String(), updated.Type)
				assert.Equal(t, "super reserve 1", updated.Description)
				assert.Equal(t, updated.Timestamp.Time.Unix(), tts.Unix())
			},
		},
		{
			Msg:      "update all information of a reserve address",
			Endpoint: fmt.Sprintf("/addresses/%d", id+99),
			Method:   http.MethodPut,
			Body: []byte(`{
"address": "0x6585Bc19A6E249acC55cd4BC38472D346a0fFaE9",
"type": "pricing_operator",
"description": "pricing operator 2"
}`),
			Assert: httputil.AssertCode(http.StatusNotFound),
		},
	}

	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}
}

func TestMain(m *testing.M) {
	var err error
	tts = time.Now().UTC()
	resolv := blockchain.NewMockContractTimestampResolver(tts)

	sugar := testutil.MustNewDevelopmentSugaredLogger()

	_, tdb = testutil.MustNewDevelopmentDB()

	tst, err = postgresql.NewStorage(sugar, tdb, resolv)
	if err != nil {
		log.Fatal(err)
	}

	ts = NewServer(sugar, "", tst)
	ts.register()

	ret := m.Run()

	if err = tst.DeleteAllTables(); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}