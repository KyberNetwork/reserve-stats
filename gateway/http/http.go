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
func NewServer(addr string,
	auth *httpsign.Authenticator,
	perm gin.HandlerFunc,
	logger *zap.Logger,
	options ...Option,
) (*Server, error) {
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

	server := Server{
		addr: addr,
		r:    r,
	}

	for _, opt := range options {
		if err := opt(&server); err != nil {
			return nil, err
		}
	}
	return &server, nil
}

// Start runs the HTTP gateway server.
func (svr *Server) Start() error {
	return svr.r.Run(svr.addr)
}
