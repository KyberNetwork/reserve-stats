package http

import (
	"net/http"
	"time"

	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is the engine to serve reserve-rate API query
type Server struct {
	r     *gin.Engine
	db    storage.ReserveRatesStorage
	host  string
	sugar *zap.SugaredLogger
}

type reserveRatesQuery struct {
	From         uint64   `form:"from" `
	To           uint64   `form:"to"`
	ReserveAddrs []string `form:"reserve" binding:"dive,isAddress"`
}

func (sv *Server) reserveRates(c *gin.Context) {
	var (
		query    reserveRatesQuery
		logger   = sv.sugar.With("func", "reserverates/http/Server.reserveRates")
		rsvAddrs []ethereum.Address
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	now := time.Now().UTC()
	if query.To == 0 {
		query.To = timeutil.TimeToTimestampMs(now)
		logger.Debug("using default to query time", "to", query.To)

		if query.From == 0 {
			query.From = timeutil.TimeToTimestampMs(now.Add(-time.Hour))
			logger = logger.With("from", query.From)
			logger.Debug("using default from query time", "from", query.From)
		}
	}

	logger = logger.With("to", query.To, "from", query.From)
	logger.Debug("querying reserve rates from database")
	for _, rsvAddr := range query.ReserveAddrs {
		rsvAddrs = append(rsvAddrs, ethereum.HexToAddress(rsvAddr))
	}
	result, err := sv.db.GetRatesByTimePoint(rsvAddrs, query.From, query.To)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "query", query)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	if result == nil {
		result = make(map[string]map[uint64]common.ReserveRates)
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
func NewServer(host string, db storage.ReserveRatesStorage, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:     r,
		db:    db,
		host:  host,
		sugar: sugar,
	}, nil
}
