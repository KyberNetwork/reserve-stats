package http

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Server is HTTP server of gateway service.
type Server struct {
	r    *gin.Engine
	addr string
}

func newReverseProxyMW(target string) (gin.HandlerFunc, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}

// NewServer creates new instance of gateway HTTP server.
func NewServer(addr, tradeLogsURL string) (*Server, error) {
	r := gin.Default()

	// TODO: use httpsignatures middleware here

	tradeLogsProxyMW, err := newReverseProxyMW(tradeLogsURL)
	if err != nil {
		return nil, err
	}

	//TODO: maps /trade-logs URLs
	r.GET("/trade-logs", tradeLogsProxyMW)

	// TODO: creates reserve-rates API and map URLs

	// TODO: creates users and map URLs

	return &Server{
		addr: addr,
		r:    r,
	}, nil
}

// Start runs the HTTP gateway server.
func (svr *Server) Start() error {
	return svr.r.Run(svr.addr)
}
