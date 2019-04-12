package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	maxTimeFrame     = time.Hour * 24 * 365 * 1 // 1 year
	defaultTimeFrame = time.Hour * 24           // 1 day
)

// Server is the HTTP server of accounting wallet-erc20-txs HTTP API.
type Server struct {
	sugar *zap.SugaredLogger
	r     *gin.Engine
	host  string
	st    storage.ReserveTransactionStorage
}

type getTransactionsQuery struct {
	httputil.TimeRangeQuery
	Wallet string `form:"wallet" binding:"isAddress"`
	Token  string `form:"token" binding:"isAddress"`
}

func (s *Server) getTransactions(c *gin.Context) {
	var (
		query getTransactionsQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime, toTime, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)

	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	data, err := s.st.GetWalletERC20Transfers(
		ethereum.HexToAddress(query.Wallet),
		ethereum.HexToAddress(query.Token),
		fromTime,
		toTime,
	)

	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, data)
}

// NewServer creates a new instance of Server.
func NewServer(sugar *zap.SugaredLogger, host string, st storage.ReserveTransactionStorage) *Server {
	r := gin.Default()
	return &Server{sugar: sugar, r: r, host: host, st: st}
}

func (s *Server) register() {
	s.r.GET("/wallet/transactions", s.getTransactions)
}

// Run starts the HTTP server and runs in foreground until terminate by user.
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
