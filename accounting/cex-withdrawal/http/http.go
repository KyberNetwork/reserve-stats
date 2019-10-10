package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	huobiStorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	maxTimeFrame     = time.Hour * 24 * 30 // 30 days
	defaultTimeFrame = time.Hour * 24      // 1 day
)

// Server is the engine to serve cex-trade-withdrawal API query
type Server struct {
	r         *gin.Engine
	huobiDB   huobiStorage.Interface
	binanceDB withdrawalstorage.Interface
	host      string
	sugar     *zap.SugaredLogger
}

type queryInput struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex" binding:"dive,isValidCEXName"`
}

type response struct {
	Huobi   []huobi.WithdrawHistory   `json:"huobi,omitempty"`
	Binance []binance.WithdrawHistory `json:"binance,omitempty"`
}

func (sv *Server) get(c *gin.Context) {
	var (
		query              queryInput
		logger             = sv.sugar.With("func", caller.GetCurrentFunctionName())
		huobiWithdrawals   []huobi.WithdrawHistory
		binanceWithdrawals []binance.WithdrawHistory
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
			common.Huobi.String(),
			common.Binance.String()}
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
		switch cex {
		case common.Huobi.String():
			huobiWithdrawals, err = sv.huobiDB.GetWithdrawHistory(from, to)
			if err != nil {
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
		case common.Binance.String():
			binanceWithdrawals, err = sv.binanceDB.GetWithdrawHistory(from, to)
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
		Huobi:   huobiWithdrawals,
		Binance: binanceWithdrawals,
	})
}

func (sv *Server) register() {
	sv.r.GET("/withdrawals", sv.get)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, huobiDB huobiStorage.Interface, binanceDB withdrawalstorage.Interface, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:         r,
		huobiDB:   huobiDB,
		binanceDB: binanceDB,
		host:      host,
		sugar:     sugar,
	}, nil
}
