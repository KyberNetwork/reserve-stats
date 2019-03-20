package listedtokenserver

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//Server struct for listed token api
type Server struct {
	sugar *zap.SugaredLogger
	r     *gin.Engine
	host  string
	//TODO: wait for
	// storage listedtokenstorage.Interface
}

type reserveTokenQuery struct {
	Reserve string `json:"reserve" binding:"required,isAddress"`
}

//NewServer return new server object
func NewServer(sugar *zap.SugaredLogger, host string) *Server {
	r := gin.Default()
	return &Server{
		sugar: sugar,
		r:     r,
		host:  host,
	}
}

func (s *Server) getReserveToken(c *gin.Context) {
	var query reserveTokenQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	//TODO: get listed token from DB
}

func (s *Server) register() {
	s.r.GET("/reserve/tokens", s.getReserveToken)
}

//Run server
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
