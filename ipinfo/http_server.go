package ipinfo

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrInvalidIP error for invalid ip input
const ErrInvalidIP = "invalid ip input"

// HTTPServer to serve endpoint
type HTTPServer struct {
	r     *gin.Engine
	l     *Locator
	host  string
	sugar *zap.SugaredLogger
}

// NewHTTPServer return an instance of HTTPServer
func NewHTTPServer(logger *zap.Logger, dataDir string, host string) (*HTTPServer, error) {
	sugar := logger.Sugar()
	l, err := NewLocator(sugar, dataDir)
	if err != nil {
		return nil, err
	}
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return &HTTPServer{
		r:     r,
		l:     l,
		host:  host,
		sugar: sugar,
	}, nil
}

func (h *HTTPServer) register() {
	h.r.GET("/ip/:ip", h.lookupIPCountry)
}

// Run start HTTPServer
func (h *HTTPServer) Run() error {
	h.register()
	return h.r.Run(h.host)
}

func (h *HTTPServer) lookupIPCountry(c *gin.Context) {
	ip := c.Param("ip")
	ipParsed := net.ParseIP(ip)
	if ipParsed == nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			errors.New(ErrInvalidIP),
		)
		return
	}
	location, err := h.l.IPToCountry(ipParsed)
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
		gin.H{
			"country": location,
		},
	)
}
