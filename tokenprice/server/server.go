package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/tokenprice/common"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage"
)

// Server serve token price via http endpoint
type Server struct {
	storage storage.Storage
	host    string
	sugar   *zap.SugaredLogger
}

// NewServer return server instance
func NewServer(sugar *zap.SugaredLogger, host string, storage storage.Storage) *Server {
	return &Server{
		storage: storage,
		host:    host,
		sugar:   sugar,
	}
}

type queryPrice struct {
	Token    string `form:"token" binding:"required"`
	Currency string `form:"currency" binding:"required"`
	Date     string `form:"date"`
}

func (sv *Server) getPrice(c *gin.Context) {
	var (
		query queryPrice
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	date := query.Date
	if len(date) == 0 {
		date = common.TimeToDateString(time.Now().UTC())
	}
	t, err := common.DateStringToTime(date)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	price, err := sv.storage.GetTokenRate(query.Token, query.Currency, t)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"price": price,
	})
}

func (sv *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/price", sv.getPrice)
	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}
