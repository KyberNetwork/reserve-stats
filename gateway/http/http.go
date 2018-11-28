package http

import (
	"net/http/httputil"
	"net/url"
	"time"

	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/httpsign"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func NewServer(addr, tradeLogsURL, reserveRatesURL, userURL, priceAnalyticURL string,
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
	if tradeLogsURL != "" {
		tradeLogsProxyMW, err := newReverseProxyMW(tradeLogsURL)
		if err != nil {
			return nil, err
		}
		r.GET("/trade-logs", tradeLogsProxyMW)
		r.GET("/burn-fee", tradeLogsProxyMW)
		r.GET("/asset-volume", tradeLogsProxyMW)
		r.GET("/reserve-volume", tradeLogsProxyMW)
		r.GET("/wallet-fee", tradeLogsProxyMW)
		r.GET("/user-volume", tradeLogsProxyMW)
		r.GET("/user-list", tradeLogsProxyMW)
		r.GET("/trade-summary", tradeLogsProxyMW)
		r.GET("/wallet-stats", tradeLogsProxyMW)
		r.GET("/country-stats", tradeLogsProxyMW)
		r.GET("/heat-map", tradeLogsProxyMW)
	}
	if reserveRatesURL != "" {
		reserveRateProxyMW, err := newReverseProxyMW(reserveRatesURL)
		if err != nil {
			return nil, err
		}
		r.GET("/reserve-rates", reserveRateProxyMW)
	}

	if userURL != "" {
		userProxyMW, err := newReverseProxyMW(userURL)
		if err != nil {
			return nil, err
		}
		r.GET("/users", userProxyMW)
		r.POST("/users", userProxyMW)
	}

	if priceAnalyticURL != "" {
		priceProxyMW, err := newReverseProxyMW(priceAnalyticURL)
		if err != nil {
			return nil, err
		}
		r.GET("/price-analytic-data", priceProxyMW)
		r.POST("/price-analytic-data", priceProxyMW)
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
