package client

import (
	"log"
	"os"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/stretchr/testify/require"
)

const testURL = "http://127.0.0.1:8009"

var cl *Client

func TestGetAllReserveAddress(t *testing.T) {
	//skip as it need external resource
	t.Skip()
	result, err := cl.GetAllReserveAddress()
	require.NoError(t, err)
	log.Printf("result is %v", result)
}

func TestMain(m *testing.M) {
	var err error
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	cl, err = NewClient(sugar, testURL)
	if err != nil {
		log.Fatal(err)
	}

	ret := m.Run()

	os.Exit(ret)
}
