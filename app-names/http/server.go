package http

import (
	"net/http"
	"strconv"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
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
	name := c.Query("name")
	logger.Debug(name)
	apps, err := sv.db.GetAllApp(name)
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
		apps,
	)
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

func (sv *Server) createApp(c *gin.Context) {
	var (
		logger   = sv.sugar.With("func", "app-names/server.createApp")
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

func (sv *Server) updateApp(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "app-names/server.updateApp")
		q     common.AppObject 
	)
	logger.Debug("start update app")
	appID, err := strconv.ParseInt(c.Param("appID"), 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if err := c.ShouldBindJSON(&q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if err := sv.db.UpdateAppAddress(appID, q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (sv *Server) deleteApp(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "app-names/server.deleteApp")
	)
	logger.Debug("delete app")
	appID, err := strconv.ParseInt(c.Param("appID"), 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if err := sv.db.DeleteApp(appID); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (sv *Server) register() {
	sv.r.GET("/application-names", sv.getApps)
	sv.r.GET("/application-names/:appID", sv.getAddressFromAppID)
	sv.r.POST("/application-names", sv.createApp)
	sv.r.PUT("/application-names/:appID", sv.updateApp)
	sv.r.DELETE("/application-names/:appID", sv.deleteApp)
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
