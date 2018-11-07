package http

import (
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//Server for price analytic service
type Server struct {
	sugar   *zap.SugaredLogger
	r       *gin.Engine
	host    string
	storage storage.Interface
}

//NewHTTPServer return new server instance
func NewHTTPServer(sugar *zap.SugaredLogger, host string, storage storage.Interface) *Server {
	r := gin.Default()
	return &Server{
		sugar:   sugar,
		r:       r,
		host:    host,
		storage: storage,
	}
}

func (s *Server) updatePriceAnalytic(c *gin.Context) {
	var priceAnalytic common.PriceAnalytic
	if err := c.ShouldBindJSON(&priceAnalytic); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if err := s.storage.UpdatePriceAnalytic(priceAnalytic); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

func (s *Server) getPriceAnalytic(c *gin.Context) {
	const maxTimeFrame = time.Hour * 24 * 365 * 3 // 3 years
	var query httputil.TimeRangeQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime, toTime, err := query.Validate(httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame))
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	priceAnalytic, err := s.storage.GetPriceAnalytic(fromTime, toTime)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(
		http.StatusOK,
		priceAnalytic,
	)
}

func (s *Server) register() {
	s.r.POST("/price-analytics", s.updatePriceAnalytic)
	s.r.GET("/price-analytics", s.getPriceAnalytic)
}

// Run server
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
