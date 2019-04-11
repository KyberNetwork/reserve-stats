package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
	"github.com/KyberNetwork/reserve-stats/app-names/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
)

// Server is the engine to serve reserve-rate API query
type Server struct {
	r     *gin.Engine
	host  string
	sugar *zap.SugaredLogger
	db    storage.Interface
}

func (sv *Server) getApps(c *gin.Context) {
	var (
		logger  = sv.sugar.With("func", "app-names/http/Server.getAddrToAppName")
		filters []storage.Filter
	)

	logger.Debug("getting addr to App name")
	name, ok := c.GetQuery("name")
	if ok {
		logger.Debugw("got name parameter from query", "name", name)
		filters = append(filters, storage.WithNameFilter(name))
	}

	activeStr, ok := c.GetQuery("active")
	if ok {
		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusBadRequest,
				err,
			)
			return
		}
		logger.Debugw("got active parameter from query", "active", active)
		if active {
			filters = append(filters, storage.WithActiveFilter())
		} else {
			filters = append(filters, storage.WithInactiveFilter())
		}
	}

	apps, err := sv.db.GetAll(filters...)
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

	result, err := sv.db.Get(appID)
	if err != nil {
		if err == storage.ErrNotExists {
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
		update bool
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
	if id, update, err = sv.db.CreateOrUpdate(q); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	app, err := sv.db.Get(id)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
	}
	if update {
		c.JSON(http.StatusOK, app)
	} else {
		c.JSON(http.StatusCreated, app)
	}
}

func (sv *Server) updateApp(c *gin.Context) {
	var q common.Application
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
	q.ID = appID

	err = sv.db.Update(q)
	if err != nil {
		if err == storage.ErrNotExists {
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
	app, err := sv.db.Get(appID)
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
	if err := sv.db.Delete(appID); err != nil {
		if err == storage.ErrNotExists {
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
func NewServer(host string, appNameDB storage.Interface, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:     r,
		db:    appNameDB,
		host:  host,
		sugar: sugar,
	}, nil
}
