package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	huobistorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage"
	"github.com/KyberNetwork/reserve-stats/accounting/zerox/storage"
)

// Server is the HTTP server of accounting CEX getTrades HTTP API.
type Server struct {
	sugar *zap.SugaredLogger
	r     *gin.Engine
	host  string
	hs    huobistorage.Interface
	bs    tradestorage.Interface
	zs    *storage.ZeroxStorage
}

// NewServer creates a new instance of Server.
func NewServer(sugar *zap.SugaredLogger, host string, hs huobistorage.Interface, bs tradestorage.Interface, zs *storage.ZeroxStorage) *Server {
	r := gin.Default()
	return &Server{
		sugar: sugar,
		r:     r,
		host:  host,
		hs:    hs,
		bs:    bs,
		zs:    zs,
	}

}

func (s *Server) register() {
	s.r.GET("/trades", s.getTrades)
	s.r.GET("/convert_to_eth_price", s.getConvertToETHPrice)
	s.r.GET("/convert_0x_trades", s.get0xConvertTrades)
	s.r.GET("/convert_trades", s.getConvertTrades)
}

// Run starts the HTTP server and runs in foreground until terminate by user.
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
