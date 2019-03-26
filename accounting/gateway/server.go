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
func NewServer(addr, listedTokenURL string,
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
	if listedTokenURL != "" {
		listedTokenProxyMW, err := newReverseProxyMW(listedTokenURL)
		if err != nil {
			return nil, err
		}
		r.GET("/reserve/tokens", listedTokenProxyMW)
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
