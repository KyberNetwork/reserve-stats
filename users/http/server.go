package http

import (
	"net/http"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, rateProvider tokenrate.ETHUSDRateProvider,
	storage storage.Interface, host string,
	influxStorage *storage.InfluxStorage) *Server {
	r := gin.Default()
	return &Server{
		sugar:         sugar,
		rateProvider:  httputil.NewCachedRateProvider(sugar, rateProvider, time.Hour),
		storage:       storage,
		r:             r,
		host:          host,
		influxStorage: influxStorage,
	}
}

//Server struct to represent a http server service
type Server struct {
	sugar         *zap.SugaredLogger
	r             *gin.Engine
	host          string
	rateProvider  tokenrate.ETHUSDRateProvider
	storage       storage.Interface
	influxStorage *storage.InfluxStorage
}

//createOrUpdate update info of an user
func (s *Server) createOrUpdate(c *gin.Context) {
	var userData common.UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if err := s.storage.CreateOrUpdate(userData); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": userData.Email})
}

func (s *Server) register() {
	s.r.POST("/users", s.createOrUpdate)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
