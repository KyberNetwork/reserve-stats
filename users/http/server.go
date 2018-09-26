package http

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Server struct to represent a http server service
type Server struct {
	r    *gin.Engine
	host string
}

//GetUsers return user info
func (s *Server) GetUsers(c *gin.Context) {

}

//GetUserInfo return infomation of an user
func (s *Server) GetUserInfo(c *gin.Context) {

}

//UpdateUserInfo update info of an user
func (s *Server) UpdateUserInfo(c *gin.Context) {

}

func (s *Server) register() {
	s.r.GET("/users", s.GetUsers)
	s.r.GET("/users/:address", s.GetUserInfo)
	s.r.POST("/users", s.UpdateUserInfo)
}

//Run start server and serve
func (s *Server) Run() {
	s.register()
	if err := s.r.Run(s.host); err != nil {
		log.Panic(err)
	}
}

//NewServer return new server instance
func NewServer(host string) *Server {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	// corsConfig.AddAllowHeaders("signed")
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &Server{r, host}
}
