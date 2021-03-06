package http

import (
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
}

// NewServer creates a new instance of Server from given parameters.
func NewServer(sugar *zap.SugaredLogger, host string, storage storage.Interface) *Server {
	r := gin.Default()
	return &Server{sugar: sugar, r: r, host: host, storage: storage}
}

func (s *Server) register() {
	s.r.POST("/addresses", s.create)
	s.r.GET("/addresses/:id", s.get)
	s.r.GET("/addresses", s.getAll)
	s.r.PUT("/addresses/:id", s.update)
}

// Run starts the HTTP server and runs in foreground until terminate by user.
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
