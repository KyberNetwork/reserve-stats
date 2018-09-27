package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

const (
	kycedTxLimit    = 6000
	nonKycedTxLimit = 3000
)

//Server struct to represent a http server service
type Server struct {
	db   *storage.UserDB
	r    *gin.Engine
	host string
}

//GetUserInfo return infomation of an user
func (s *Server) GetUserInfo(c *gin.Context) {
	address := c.Param("address")
	_, err := s.db.GetUserInfo(address)
	if err != nil {
		if err != pg.ErrNoRows {
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
		response := common.UserResponse{
			KYC: false,
			Cap: nonKycedTxLimit,
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": true,
				"data":    response,
			},
		)
		return
	}
	response := common.UserResponse{
		KYC: true,
		Cap: kycedTxLimit,
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    response,
		},
	)
}

//UpdateUserInfo update info of an user
func (s *Server) UpdateUserInfo(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "Request malformed",
			},
		)
		return
	}
	postForm := c.Request.Form
	email := postForm.Get("user")
	zap.S().Infof("email: %s", email)

	addresses := postForm.Get("addresses")
	zap.S().Infof("addresses: %s", addresses)

	times := postForm.Get("timestamps")
	zap.S().Infof("timestamps: %s", times)

	addrsStr := strings.Split(addresses, "-")
	timesStr := strings.Split(times, "-")
	if len(addrsStr) != len(timesStr) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "",
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
	if len(userAddresses) == 0 {
		return
	}

	if err := s.db.StoreUserInfo(email, userAddresses); err != nil {
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
func NewServer(userDB *storage.UserDB, host string) *Server {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	// corsConfig.AddAllowHeaders("signed")
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &Server{userDB, r, host}
}
