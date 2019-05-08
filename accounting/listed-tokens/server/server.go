package server

import (
	"log"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
)

//Server struct for listed token api
type Server struct {
	sugar   *zap.SugaredLogger
	r       *gin.Engine
	host    string
	storage storage.Interface
}

type reserveTokenQuery struct {
	Reserve string `form:"reserve" binding:"isAddress"`
}

//NewServer return new server object
func NewServer(logger *zap.Logger, host string, storage storage.Interface) *Server {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	sugar := logger.Sugar()
	return &Server{
		sugar:   sugar,
		r:       r,
		host:    host,
		storage: storage,
	}
}

func (s *Server) getReserveToken(c *gin.Context) {
	var (
		query reserveTokenQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	log.Printf("reserve: %s", query.Reserve)
	listedTokens, version, blockNumber, err := s.storage.GetTokens(ethereum.HexToAddress(query.Reserve))
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"version":      version,
			"block_number": blockNumber,
			"data":         listedTokens,
		},
	)
}

func (s *Server) register() {
	s.r.GET("/reserve/tokens", s.getReserveToken)
}

//Run server
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
