package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/zerox/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

// Server is the HTTP server of accounting CEX getTrades HTTP API.
type Server struct {
	sugar *zap.SugaredLogger
	r     *gin.Engine
	host  string
	zs    *storage.ZeroxStorage
}

// NewServer creates a new instance of Server.
func NewServer(sugar *zap.SugaredLogger, host string, zs *storage.ZeroxStorage) *Server {
	r := gin.Default()
	return &Server{
		sugar: sugar,
		r:     r,
		host:  host,
		zs:    zs,
	}

}

func (s *Server) register() {
	// s.r.GET("/trades", s.getTrades)
	s.r.GET("/convert_0x_trade", s.getConvert0xTrade)
}

// Run starts the HTTP server and runs in foreground until terminate by user.
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}

type getConvertTradeQuery struct {
	httputil.TimeRangeQuery
}

func (s *Server) getConvert0xTrade(c *gin.Context) {
	var (
		query getConvertTradeQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	result, err := s.zs.GetConvertTrades(int64(query.From), int64(query.To))
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
		result,
	)
}
