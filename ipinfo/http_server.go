package ipinfo

import (
	"fmt"
	"net/http"

	"github.com/KyberNetwork/reserve-stats/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTPServer to serve endpoint
type HTTPServer struct {
	r     *gin.Engine
	l     *Locator
	sugar *zap.SugaredLogger
}

// NewHTTPServer return an instance of HTTPServer
func NewHTTPServer(sugar *zap.SugaredLogger, dataDir string) (*HTTPServer, error) {
	l, err := NewLocator(sugar, dataDir)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{
		r:     gin.Default(),
		l:     l,
		sugar: sugar,
	}, nil
}

// Run start HTTPServer
func (h *HTTPServer) Run() error {
	port := fmt.Sprintf(":%d", util.IPLocatorPort)
	h.r.GET("/ip/:ip", h.checkIPLocator)
	return h.r.Run(port)
}

func (h *HTTPServer) checkIPLocator(c *gin.Context) {
	ip := c.Param("ip")
	location, err := h.l.IPToCountry(ip)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
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
