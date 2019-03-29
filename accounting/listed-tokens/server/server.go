package server

import (
	"net/http"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
)

//Server struct for listed token api
type Server struct {
	sugar   *zap.SugaredLogger
	r       *gin.Engine
	host    string
	storage storage.Interface
}

type reserveTokenQuery struct {
	Reserve string `form:"reserve" binding:"required,isAddress"`
}

//NewServer return new server object
func NewServer(sugar *zap.SugaredLogger, host string, storage storage.Interface) *Server {
	r := gin.Default()
	return &Server{
		sugar:   sugar,
		r:       r,
		host:    host,
		storage: storage,
	}
}

func (s *Server) getReserveToken(c *gin.Context) {
	var query reserveTokenQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	reserveAddr := ethereum.HexToAddress(query.Reserve)

	listedTokens, version, blockNumber, err := s.storage.GetTokens(reserveAddr.Hex())
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
