package http

import (
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage"
)

// Server is the HTTP server of accounting reserve addresses service.
type Server struct {
	sugar   *zap.SugaredLogger
	r       *gin.Engine
	host    string
	storage storage.Interface
	resolv  blockchain.ContractTimestampResolver
}

// NewServer creates a new instance of Server from given parameters.
func NewServer(sugar *zap.SugaredLogger, host string, storage storage.Interface, resolv blockchain.ContractTimestampResolver) *Server {
	r := gin.Default()
	return &Server{sugar: sugar, r: r, host: host, storage: storage, resolv: resolv}
}

func (s *Server) register() {
	s.r.POST("/addresses", s.create)
	s.r.GET("/addresses/:id", s.get)
	s.r.PUT("/addresses/:id", nil)
}

func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
