package http

import (
	"net/http"

	appNames "github.com/KyberNetwork/reserve-stats/integration-app-names"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is the engine to serve reserve-rate API query
type Server struct {
	r             *gin.Engine
	host          string
	sugar         *zap.SugaredLogger
	addrToAppName appNames.AddrToAppName
}

// AddrNameQuery define a struct to contain param for a POST add-to-appname
type AddrNameQuery struct {
	Address string `json:"address" binding:"required,isAddress"`
	AppName string `json:"appname" binding:"required"`
}

func (sv *Server) getAddrToAppName(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "intergration-app-names/http/Server.getAddrToAppName")
	)

	logger.Debug("getting addr to App name")
	result := sv.addrToAppName.GetAddrToAppName()
	c.JSON(http.StatusOK, result)
}

func (sv *Server) updateAddrToAppName(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "intergration-app-names/http/Server.updateAddrToAppName")
		q      AddrNameQuery
	)

	logger.Debug("updating addr to App name")
	if err := c.ShouldBindJSON(&q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	if err := sv.addrToAppName.UpdateMapAddrAppName(ethereum.HexToAddress(q.Address), q.AppName); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{q.Address: q.AppName})
}

func (sv *Server) register() {
	sv.r.GET("/addr-to-appname", sv.getAddrToAppName)
	sv.r.POST("/addr-to-appname", sv.updateAddrToAppName)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, atan appNames.AddrToAppName, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:             r,
		addrToAppName: atan,
		host:          host,
		sugar:         sugar,
	}, nil
}
