package http

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/KyberNetwork/reserve-stats/users/cmc"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"go.uber.org/zap"
)

func connectToTestDB() *storage.UserDB {
	return storage.NewDB(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRESS_DATABASE"),
	)
}

func clearTestDB() {

}

func TestUserHTTPServer(t *testing.T) {
	userStorage := connectToTestDB()
	cmc := cmc.NewCMCEthUSDRate()
	userStats := stats.NewUserStats(cmc, userStorage)
	defer clearTestDB()

	// create new Server instance
	//TODO: turn into variable
	host := fmt.Sprintf(":%s", "9000")
	s := NewServer(userStats, host)

	zap.S().Infof("Server instance: %+v", s)

	// test case
	const (
		requestEndpoint = "/users"
		userEmail       = "test@gmail.com"
		wrongUserEmail  = "test"

		userAddresses          = "0xc9a658f87d7432ff897f31dce318f0856f66acb7_0x2ea6200a999f4c6c982be525f8dc294f14f4cb08"
		wrongUserAddresses     = "wrong-address_0x13197"
		wrongNumberOfAddresses = "0xc9a658f87d7432ff897f31dce318f0856f66acb7_0x2ea6200a999f4c6c982be525f8dc294f14f4cb08_0x4e012a6445ba2a590b8b1ee4e95d03e345a0c2e5"
		userTimeStamp          = "1538380670000_1538380682000"
	)

	var tests = []testCase{
		{
			msg:      "test get empty db",
			endpoint: requestEndpoint,
			method:   http.MethodGet,
			assertFn: httputil.ExpectSuccess,
		},
	}

}
