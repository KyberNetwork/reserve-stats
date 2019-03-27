package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
)

//default frame is 1 year =365  days
var maxTimeFrame = 365 * 24 * time.Hour

// Server is the engine to serve reserve-rate API query
type Server struct {
	r     *gin.Engine
	db    storage.Interface
	host  string
	sugar *zap.SugaredLogger
}

func (sv *Server) reserveRates(c *gin.Context) {
	var (
		query  httputil.TimeRangeQuery
		logger = sv.sugar.With("func", "accouting/reserve-rate/http/Server.reserveRates")
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	from, to, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
	)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	logger = logger.With("from", from, "to", to)
	logger.Debug("querying rates from database")

	reserveRate, err := sv.db.GetRates(from, to)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	ethUsdRate, err := sv.db.GetETHUSDRates(from, to)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	result := storage.AccountingRatesReply{
		ReserveRates: reserveRate,
		EthUsdRates:  ethUsdRate,
	}

	c.JSON(http.StatusOK, result)
}

func (sv *Server) register() {
	sv.r.GET("/reserve-rates", sv.reserveRates)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, db storage.Interface, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:     r,
		db:    db,
		host:  host,
		sugar: sugar,
	}, nil
}
