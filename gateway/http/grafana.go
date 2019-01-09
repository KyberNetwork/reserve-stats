package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	grafanaKeyID   = "Bearer"
	grafanaPrexfix = "grafana"
)

//Since the gateway group grafana request by /grafana/..., this grafana/ must be remove
func removeFirstComponent(path string) string {
	return strings.TrimPrefix(path, "/"+grafanaPrexfix)

}

//newGrafanaDirector return a director function to be used in GrafanaProxy
func newGrafanaDirector(targetURL string, apiKey string) (func(req *http.Request), error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, removeFirstComponent(req.URL.Path))
		log.Print("redirecting...")
		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		log.Printf("%s", req.URL)
		req.Header.Add("Authorization", grafanaKeyID+" "+apiKey)
		log.Printf("%s", req.Header.Get("Authorization"))
		log.Printf("%s", req.Header)
	}, nil
}

//NewGrafanaProxy return a proxy for grafana
func NewGrafanaProxy(grafanaURL string, apiKey string) (gin.HandlerFunc, error) {
	grafanaDirector, err := newGrafanaDirector(grafanaURL, apiKey)
	if err != nil {
		return nil, err
	}
	proxy := &httputil.ReverseProxy{
		Director: grafanaDirector,
	}
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}, nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
