package http

import (
	"net/http/httputil"
	"net/url"
	"time"

	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-contrib/cors"
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
func NewServer(addr, tradeLogsURL, reserveRatesURL, userURL, priceAnalyticURL, readKeyID, readKeySecret, writeKeyID, writeKeySecret string) (*Server, error) {
	r := gin.Default()
	r.Use(libhttputil.MiddlewareHandler)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Digest", "Authorization", "Signature", "Nonce")
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	// signature middleware for signing message
	auth, err := newAuthenticator(readKeyID, readKeySecret, writeKeyID, writeKeySecret)
	// Permision middleware for checking permission
	perm, err := newPermissioner(readKeyID, writeKeyID)
	if err != nil {
		return nil, err
	}
	if tradeLogsURL != "" {
		tradeLogsProxyMW, err := newReverseProxyMW(tradeLogsURL)
		if err != nil {
			return nil, err
		}
		authGroup := r.Group("/")
		authGroup.Use(perm)
		authGroup.Use(auth.Authenticated())
		authGroup.GET("/trade-logs", tradeLogsProxyMW)
		authGroup.GET("/burn-fee", tradeLogsProxyMW)
		authGroup.GET("/asset-volume", tradeLogsProxyMW)
		authGroup.GET("/reserve-volume", tradeLogsProxyMW)
		authGroup.GET("/wallet-fee", tradeLogsProxyMW)
		authGroup.GET("/user-volume", tradeLogsProxyMW)
		authGroup.GET("/user-list", tradeLogsProxyMW)
		authGroup.GET("/trade-summary", tradeLogsProxyMW)
		authGroup.GET("/wallet-stats", tradeLogsProxyMW)
		authGroup.GET("/country-stats", tradeLogsProxyMW)
		authGroup.GET("/heat-map", tradeLogsProxyMW)
	}
	if reserveRatesURL != "" {
		reserveRateProxyMW, err := newReverseProxyMW(reserveRatesURL)
		if err != nil {
			return nil, err
		}
		authGroup := r.Group("/")
		authGroup.Use(perm)
		authGroup.Use(auth.Authenticated())
		authGroup.GET("/reserve-rates", reserveRateProxyMW)
	}

	if userURL != "" {
		userProxyMW, err := newReverseProxyMW(userURL)
		if err != nil {
			return nil, err
		}
		authGroup := r.Group("/")
		authGroup.Use(perm)
		authGroup.Use(auth.Authenticated())
		authGroup.GET("/users", userProxyMW)
		authGroup.POST("/users", userProxyMW)
	}

	if priceAnalyticURL != "" {
		priceProxyMW, err := newReverseProxyMW(priceAnalyticURL)
		if err != nil {
			return nil, err
		}
		authGroup := r.Group("/")
		authGroup.Use(perm)
		authGroup.Use(auth.Authenticated())
		authGroup.GET("/price-analytic-data", priceProxyMW)
		authGroup.POST("/price-analytic-data", priceProxyMW)
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
