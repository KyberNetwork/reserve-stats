package http

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/stats"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v8"
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

func isAddress(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	address := field.String()
	if err := validation.Validate(address, validation.Required); err != nil {
		return false
	}
	a := ethereum.HexToAddress(address)
	if a.Big().Cmp(ethereum.Big0) == 0 {
		return false
	}
	return true
}

//IsEmail validation function for email field
func IsEmail(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	if err := validation.Validate(field.String(), is.Email); err != nil {
		return false
	}
	return true
}

//UpdateUserInfo update info of an user
func (s *Server) UpdateUserInfo(c *gin.Context) {
	var userData common.UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	err := s.userStats.StoreUserInfo(userData)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isAddress", isAddress)
		v.RegisterValidation("isemail", IsEmail)
	}
	return &Server{sugar, userStats, r, host}
}
