package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/depositstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	maxTimeFrame     = time.Hour * 24 * 30 // 30 days
	defaultTimeFrame = time.Hour * 24      // 1 day
)

// Server is the engine to serve cex-trade-withdrawal API query
type Server struct {
	r         *gin.Engine
	binanceDB *depositstorage.BinanceStorage
	host      string
	sugar     *zap.SugaredLogger
}

type queryInput struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex"`
}

type response struct {
	Binance map[string][]binance.DepositHistory `json:"binance,omitempty"`
}

func (sv *Server) get(c *gin.Context) {
	var (
		query  queryInput
		logger = sv.sugar.With("func", caller.GetCurrentFunctionName())
		// huobiWithdrawals = make(map[string][]huobi.WithdrawHistory)
		binanceDeposit = make(map[string][]binance.DepositHistory) // map account with its trades
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if len(query.Exchanges) == 0 {
		query.Exchanges = []string{
			common.Binance.String(),
		}
	}

	from, to, err := query.Validate(
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

	logger = logger.With("from", from, "to", to, "exchanges", query.Exchanges)
	logger.Debug("querying withdrawals from database")

	for _, cex := range query.Exchanges {
		if cex == common.Binance.String() {
			binanceDeposit, err = sv.binanceDB.GetDepositHistory(from, to)
			if err != nil {
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
		}
	}

	c.JSON(http.StatusOK, response{
		Binance: binanceDeposit,
	})
}

func (sv *Server) register() {
	sv.r.GET("/deposits", sv.get)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, binanceDB *depositstorage.BinanceStorage, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:         r,
		binanceDB: binanceDB,
		host:      host,
		sugar:     sugar,
	}, nil
}
