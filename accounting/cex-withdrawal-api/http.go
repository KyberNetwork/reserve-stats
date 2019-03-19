package cexwithdrawalapi

import (
	"net/http"
	"time"

	binanceStorage "github.com/KyberNetwork/reserve-stats/accounting/binance-storage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	huobiStorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//default frame is 1 year =365  days
var maxTimeFrame = 365 * 24 * time.Hour

// Server is the engine to serve cex-withdrawal API query
type Server struct {
	r         *gin.Engine
	huobidb   huobiStorage.Interface
	binancedb binanceStorage.Interface
	host      string
	sugar     *zap.SugaredLogger
}

type queryInput struct {
	httputil.TimeRangeQuery
	Cex string `form:"cex" binding:"required,isValidCexName"`
}

func (sv *Server) get(c *gin.Context) {
	var (
		query  queryInput
		logger = sv.sugar.With("func", "accouting/cex-withdrawal-api/http.get")
		result interface{}
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	from, to, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	logger = logger.With("from", from, "to", to, "cex name", query.Cex)
	logger.Debug("querying rates from database")

	switch query.Cex {
	case common.Huobi.String():
		result, err = sv.huobidb.GetWithdrawHistory(from, to)
	case common.Binance.String():
		result, err = sv.binancedb.GetWithdrawHistory(from, to)
	default:
		//should not get here but it's a safe guard
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "exchange name is invalid"},
		)
		return
	}

	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, result)
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
func NewServer(host string, huobidb huobiStorage.Interface, binancedb binanceStorage.Interface, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:         r,
		huobidb:   huobidb,
		binancedb: binancedb,
		host:      host,
		sugar:     sugar,
	}, nil
}
