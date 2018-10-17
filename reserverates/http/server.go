package http

import (
	"net/http"

	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
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
	From         uint64   `form:"from" binding:"required"`
	To           uint64   `form:"to" binding:"required"`
	ReserveAddrs []string `form:"reserve" binding:"dive,isAddress"`
}

func (sv *Server) reserveRates(c *gin.Context) {
	var query reserveRatesQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	result, err := sv.db.GetRatesByTimePoint(query.ReserveAddrs, query.From, query.To)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
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
