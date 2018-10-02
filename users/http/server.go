package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//Server struct to represent a http server service
type Server struct {
	userStats *stats.UserStats
	r         *gin.Engine
	host      string
}

//GetUserInfo return infomation of an user
func (s *Server) GetUserInfo(c *gin.Context) {
	address := c.Param("address")
	txLimit, kyced, err := s.userStats.GetTxCapByAddress(address)
	if err != nil {
		zap.S().Errorf("Cannot get user info: %+v", err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  fmt.Sprintf("Cannot get user info: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    txLimit,
			"kyced":   kyced,
		},
	)
}

//UpdateUserInfo update info of an user
func (s *Server) UpdateUserInfo(c *gin.Context) {
	err := c.Request.ParseForm()
	log.Printf("Request: %+v", c.Request)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  fmt.Sprintf("Request malformed: %s", err.Error()),
			},
		)
		return
	}
	postForm := c.Request.Form
	email := postForm.Get("user")

	addresses := postForm.Get("addresses")

	times := postForm.Get("timestamps")

	addrsStr := strings.Split(addresses, "-")
	timesStr := strings.Split(times, "-")
	if len(addrsStr) != len(timesStr) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "Addresses and timestamps length does not match",
			},
		)
		return
	}
	userAddresses := []common.UserAddress{}
	for i, addr := range addrsStr {
		var (
			t uint64
			a = ethereum.HexToAddress(addr)
		)
		t, err = strconv.ParseUint(timesStr[i], 10, 64)
		if a.Big().Cmp(ethereum.Big0) != 0 && err == nil {
			userAddresses = append(userAddresses, common.UserAddress{
				Address:   addr,
				Timestamp: t,
			})
		}
	}
	log.Printf("user addresses: %v", userAddresses)
	if len(userAddresses) == 0 {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  fmt.Sprintf("userAddresses should not be empty"),
			},
		)
		return
	}

	if err := s.userStats.StoreUserInfo(email, userAddresses); err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  fmt.Sprintf("Cannot store user info: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
		},
	)
}

func (s *Server) register() {
	s.r.GET("/users/:address", s.GetUserInfo)
	s.r.POST("/users", s.UpdateUserInfo)
}

//Run start server and serve
func (s *Server) Run() {
	s.register()
	if err := s.r.Run(s.host); err != nil {
		zap.S().Panic(err)
	}
}

//NewServer return new server instance
func NewServer(userStats *stats.UserStats, host string) *Server {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	// corsConfig.AddAllowHeaders("signed")
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &Server{userStats, r, host}
}
