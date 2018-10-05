package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.uber.org/zap"
)

//Server struct to represent a http server service
type Server struct {
	sugar     *zap.SugaredLogger
	userStats *stats.UserStats
	r         *gin.Engine
	host      string
}

//GetUserInfo return infomation of an user
func (s *Server) GetUserInfo(c *gin.Context) {
	address := c.Param("address")
	txLimit, kyced, err := s.userStats.GetTxCapByAddress(address)
	if err != nil {
		s.sugar.Errorf("Cannot get user info: %+v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("Cannot get user info: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data":  txLimit,
			"kyced": kyced,
		},
	)
}

func getUserField(c *gin.Context) (string, error) {
	postForm := c.Request.Form
	email := postForm.Get("user")

	// validate email
	err := validation.Validate(email,
		validation.Required, // not empty
		is.Email,            // is a valid email
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"reason": fmt.Sprintf("Email is not valid: %s", err.Error()),
			},
		)
		return "", err
	}
	return email, nil
}

func getAddresses(c *gin.Context) ([]string, error) {
	var addresses []string
	addrs := c.Request.Form.Get("addresses")
	if err := validation.Validate(addrs, validation.Required); err != nil {
		return addresses, err
	}
	addresses = strings.Split(addrs, "-")
	for _, addr := range addresses {
		a := ethereum.HexToAddress(addr)
		if a.Big().Cmp(ethereum.Big0) == 0 {
			return addresses, fmt.Errorf("address %s is not valid", addr)
		}
	}
	return addresses, nil
}

func getTimestamps(c *gin.Context) ([]time.Time, error) {
	times := []time.Time{}
	timestamps := c.Request.Form.Get("timestamps")
	if err := validation.Validate(timestamps, validation.Required); err != nil {
		return times, err
	}
	timesArr := strings.Split(timestamps, "-")
	for _, t := range timesArr {
		timeInt, err := strconv.ParseInt(t, 10, 64)
		// convert to second
		timeInt = timeInt / 1000
		if err != nil {
			return times, err
		}
		nativeTime := time.Unix(timeInt, 0)
		times = append(times, nativeTime)
	}
	return times, nil
}

//UpdateUserInfo update info of an user
func (s *Server) UpdateUserInfo(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("Request malformed: %s", err.Error()),
			},
		)
		return
	}
	email, err := getUserField(c)
	if err != nil {
		return
	}
	addresses, err := getAddresses(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("addresses error: %s", err.Error()),
			},
		)
		return
	}

	times, err := getTimestamps(c)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("timestamps error: %s", err.Error()),
			},
		)
		return
	}

	if len(addresses) != len(times) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Addresses and timestamps length does not match",
			},
		)
		return
	}
	userAddresses := []common.UserAddress{}
	for i, addr := range addresses {
		userAddresses = append(userAddresses, common.UserAddress{
			Address:   addr,
			Timestamp: times[i],
		})
	}
	if len(userAddresses) == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("userAddresses should not be empty"),
			},
		)
		return
	}

	if err := s.userStats.StoreUserInfo(email, userAddresses); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("Cannot store user info: %s", err.Error()),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{},
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
		s.sugar.Panic(err)
	}
}

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, userStats *stats.UserStats, host string) *Server {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &Server{sugar: sugar, userStats: userStats, r: r, host: host}
}
