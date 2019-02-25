package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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
				assert.Equal(t, stored.Timestamp.Unix(), tts.Unix())
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

func TestMain(m *testing.M) {
	tts = time.Now().UTC()
	resolv := blockchain.NewMockContractTimestampResolver(tts)

	sugar := testutil.MustNewDevelopmentSugaredLogger()

	_, tdb = testutil.MustNewDevelopmentDB()

	st, err := postgresql.NewStorage(sugar, tdb)
	if err != nil {
		log.Fatal(err)
	}

	ts = NewServer(sugar, "", st, resolv)
	ts.register()

	ret := m.Run()

	if err = st.DeleteAllTables(); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}
