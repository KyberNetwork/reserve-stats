package http

import (
	"errors"
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
	logger.Debugw("got name parameter from query", "name", name)
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
	// TODO: using ShouldBindUri when gin support it in new release
	appIDStr := c.Param("id")
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	result, err := sv.db.GetApp(appID)
	if err != nil {
		if err == storage.ErrAppNotExist {
			httputil.ResponseFailure(
				c,
				http.StatusNotFound,
				err,
			)
		} else {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (sv *Server) createApp(c *gin.Context) {
	var (
		logger = sv.sugar.With("func", "app-names/server.createApp")
		q      common.Application
		err    error
		id     int64
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
	if id, err = sv.db.CreateOrUpdate(q); err != nil {
		if err == storage.ErrAddrExisted {
			httputil.ResponseFailure(
				c,
				http.StatusConflict,
				err,
			)
		} else {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
		}
		return
	}
	app, err := sv.db.GetApp(id)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
	}
	if q.ID != 0 {
		c.JSON(http.StatusOK, app)
	} else {
		c.JSON(http.StatusCreated, app)
	}
}

func (sv *Server) updateApp(c *gin.Context) {
	var (
		q common.Application
	)
	//TODO: using ShouldBindUri when gin support it in new release
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || appID == 0 {
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusBadRequest,
				err,
			)
			return
		}
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			errors.New("invalid app id"),
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
	err = sv.db.UpdateAppAddress(appID, q)
	if err != nil {
		if err == storage.ErrAppNotExist {
			httputil.ResponseFailure(
				c,
				http.StatusNotFound,
				err,
			)
		} else {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
		}
		return
	}
	app, err := sv.db.GetApp(appID)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
	}
	c.JSON(http.StatusOK, app)
}

func (sv *Server) deleteApp(c *gin.Context) {
	//TODO: using ShouldBindUri when gin support it in new release
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if err := sv.db.DeleteApp(appID); err != nil {
		if err == storage.ErrAppNotExist {
			httputil.ResponseFailure(
				c,
				http.StatusNotFound,
				err,
			)
		} else {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (sv *Server) register() {
	sv.r.GET("/applications", sv.getApps)
	sv.r.GET("/applications/:id", sv.getAddressFromAppID)
	sv.r.POST("/applications", sv.createApp)
	sv.r.PUT("/applications/:id", sv.updateApp)
	sv.r.DELETE("/applications/:id", sv.deleteApp)
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
