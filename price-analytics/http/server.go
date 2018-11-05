package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/price-analytics/common"
	"github.com/KyberNetwork/reserve-stats/price-analytics/storage"
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{},
		)
		return
	}
	if err := s.storage.UpdatePriceAnalytic(priceAnalytic); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{},
	)
}

func (s *Server) validateTimeInput(c *gin.Context) (time.Time, time.Time, bool) {
	var (
		from, to time.Time
	)
	fromTime, ok := strconv.ParseUint(c.Query("fromTime"), 10, 64)
	if ok != nil {
		return from, to, false
	}
	from = timeutil.TimestampMsToTime(fromTime)
	toTime, _ := strconv.ParseUint(c.Query("toTime"), 10, 64)
	if toTime == 0 {
		to = time.Now()
	} else {
		to = timeutil.TimestampMsToTime(toTime)
	}
	return from, to, true
}

func (s *Server) getPriceAnalytic(c *gin.Context) {
	fromTime, toTime, ok := s.validateTimeInput(c)
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "time input is not valid",
			},
		)
		return
	}
	priceAnalytic, err := s.storage.GetPriceAnalytic(fromTime, toTime)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		priceAnalytic,
	)
}

func (s *Server) register() {
	s.r.POST("/price-analytic-data", s.updatePriceAnalytic)
	s.r.GET("/price-analytic-data", s.getPriceAnalytic)
}

// Run server
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
