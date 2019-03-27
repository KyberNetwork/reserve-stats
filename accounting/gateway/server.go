package gateway

import (
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/httpsign"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
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
func NewServer(addr, cexTradeURL, reserveAddressURL string,
	auth *httpsign.Authenticator,
	perm gin.HandlerFunc,
	logger *zap.Logger) (*Server, error) {
	r := gin.Default()
	r.Use(libhttputil.MiddlewareHandler)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Digest", "Authorization", "Signature", "Nonce")
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(cors.New(corsConfig))
	r.Use(perm)
	r.Use(auth.Authenticated())
	if cexTradeURL != "" {
		cexTradeURLMW, err := newReverseProxyMW(cexTradeURL)
		if err != nil {
			return nil, err
		}
		r.GET("/cex_trades", cexTradeURLMW)
	}

	if reserveAddressURL != "" {
		reserveAddressURLMW, err := newReverseProxyMW(reserveAddressURL)
		if err != nil {
			return nil, err
		}
		r.POST("/addresses", reserveAddressURLMW)
		r.GET("/addresses/:id", reserveAddressURLMW)
		r.GET("/addresses", reserveAddressURLMW)
		r.PUT("/addresses/:id", reserveAddressURLMW)
	}

	return &Server{
		addr: addr,
		r:    r,
	}, nil
}

// Start runs the HTTP gateway server.
func (svr *Server) Start() error {
	return svr.r.Run(svr.addr)
}
