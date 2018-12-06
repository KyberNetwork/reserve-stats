package http

import (
	"net/http"
	"strconv"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	// ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is the engine to serve reserve-rate API query
type Server struct {
	r     *gin.Engine
	host  string
	sugar *zap.SugaredLogger
	db    *storage.AppNameDB
}

func (sv *Server) getApps(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "intergration-app-names/http/Server.getAddrToAppName")
	)

	logger.Debug("getting addr to App name")
}

func (sv *Server) getAddressFromAppID(c *gin.Context) {
	appIDStr := c.Param("appID")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	result, err := sv.db.GetAppAddresses(appID)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (sv *Server) updateAddrToAppName(c *gin.Context) {
	var (
		logger   = sv.sugar.With("func", "intergration-app-names/http/Server.updateAddrToAppName")
		q        common.AppObject
		response common.AppObject
		err      error
	)

	logger.Debug("updating addr to App name")
	if err := c.ShouldBindJSON(&q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if response, err = sv.db.CreateOrUpdate(q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (sv *Server) register() {
	sv.r.GET("/application-names", sv.getApps)
	sv.r.GET("/application-names/:appID", sv.getAddressFromAppID)
	sv.r.POST("/application-name", sv.updateAddrToAppName)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, appNameDB *storage.AppNameDB, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:     r,
		db:    appNameDB,
		host:  host,
		sugar: sugar,
	}, nil
}
