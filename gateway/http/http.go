package http

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/contrib/httpsignatures"
	"github.com/gin-gonic/contrib/httpsignatures/crypto"
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
func NewServer(addr, tradeLogsURL, reserveRatesURL, userURL, secretKey string) (*Server, error) {
	r := gin.Default()

	// signature middleware for signing message
	hmacsha512 := &crypto.HmacSha512{}
	signKeyID := httpsignatures.KeyID("sign")
	secrets := httpsignatures.Secrets{
		signKeyID: &httpsignatures.Secret{
			Key:       secretKey,
			Algorithm: hmacsha512,
		},
	}
	auth := httpsignatures.NewAuthenticator(secrets)

	if tradeLogsURL != "" {
		tradeLogsProxyMW, err := newReverseProxyMW(tradeLogsURL)
		if err != nil {
			return nil, err
		}
		r.GET("/trade-logs", tradeLogsProxyMW)
		//
		// r.GET("/burn-fee", tradeLogsProxyMW)
		// r.GET("/wallet-fee", tradeLogsProxyMW)
		// r.GET("/country-stats", tradeLogsProxyMW)
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
		r.POST("/users", auth.Authenticated(), userProxyMW)
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
