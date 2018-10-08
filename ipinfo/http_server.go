package ipinfo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTPServer to serve endpoint
type HTTPServer struct {
	r     *gin.Engine
	l     *Locator
	port  int
	sugar *zap.SugaredLogger
}

// NewHTTPServer return an instance of HTTPServer
func NewHTTPServer(sugar *zap.SugaredLogger, dataDir string, port int) (*HTTPServer, error) {
	l, err := NewLocator(sugar, dataDir)
	if err != nil {
		return nil, err
	}
	return &HTTPServer{
		r:     gin.Default(),
		l:     l,
		port:  port,
		sugar: sugar,
	}, nil
}

func (h *HTTPServer) register() {
	h.r.GET("/ip/:ip", h.lookupIPCountry)
}

// Run start HTTPServer
func (h *HTTPServer) Run() error {
	h.register()
	port := fmt.Sprintf(":%d", h.port)
	return h.r.Run(port)
}

func (h *HTTPServer) lookupIPCountry(c *gin.Context) {
	ip := c.Param("ip")
	location, err := h.l.IPToCountry(ip)
	if err != nil {
		responseCode := http.StatusBadRequest
		if err != ErrInvalidIP {
			responseCode = http.StatusInternalServerError
		}
		c.JSON(
			responseCode,
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
