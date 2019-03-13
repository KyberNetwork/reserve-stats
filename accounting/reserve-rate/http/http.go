package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
		logger = sv.sugar.With("func", "reserverates/http/Server.reserveRates")
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	_, _, err := query.Validate()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	logger = logger.With("to", query.To, "from", query.From)
	logger.Debug("querying reserve rates from database")

	reserveRate, err := sv.db.GetRates(timeutil.TimestampMsToTime(query.From), timeutil.TimestampMsToTime(query.To))
	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	ethUsdRate, err := sv.db.GetETHUSDRates(timeutil.TimestampMsToTime(query.From), timeutil.TimestampMsToTime(query.To))
	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
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
