package http

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/tokenrate"
	"net/http"
	"reflect"
	"time"

	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v8"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, rateProvider tokenrate.ETHUSDRateProvider, storage storage.Interface, host string) *Server {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isAddress", isAddress)
		v.RegisterValidation("isemail", IsEmail)
	}
	return &Server{
		sugar:        sugar,
		rateProvider: newCachedRateProvider(sugar, rateProvider),
		storage:      storage,
		r:            r, host: host}
}

//Server struct to represent a http server service
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	storage      storage.Interface
}

//IsKYCed return infomation of an user
func (s *Server) getTransactionLimit(c *gin.Context) {
	address := c.Param("address")

	kyced, err := s.storage.IsKYCed(address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to check KYC status: %s", err.Error()),
			},
		)
		return
	}

	uc := common.NewUserCap(kyced)
	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to get usd rate: %s", err.Error()),
			},
		)
		return
	}

	// maximum of ETH in wei
	txLimit := blockchain.EthToWei(uc.TxLimit / rate)

	c.JSON(
		http.StatusOK,
		gin.H{
			"data":  txLimit,
			"kyced": kyced,
		},
	)
}

func isAddress(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
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
func IsEmail(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	if err := validation.Validate(field.String(), is.Email); err != nil {
		return false
	}
	return true
}

//createOrUpdate update info of an user
func (s *Server) createOrUpdate(c *gin.Context) {
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

	if err := s.storage.CreateOrUpdate(userData); err != nil {
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
	s.r.GET("/users/:address", s.getTransactionLimit)
	s.r.POST("/users", s.createOrUpdate)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
