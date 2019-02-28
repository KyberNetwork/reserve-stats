package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
)

type createInput struct {
	Address     string `json:"address" binding:"required,isAddress"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (s *Server) create(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/http/*Server.create")
		input  createInput
		ts     time.Time
	)

	err := c.ShouldBindJSON(&input)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	logger = logger.With(
		"address", input.Address,
		"type", input.Type,
		"description", input.Description,
	)

	logger.Debug("received request")

	address := ethereum.HexToAddress(input.Address)

	if ts, err = s.resolv.Resolve(address); err == blockchain.ErrNotAvailble {
		logger.Infow("address creation time is not available", "err", err.Error())
		err = nil
	} else if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	addressType, ok := common.IsValidAddressType(input.Type)
	if !ok {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			fmt.Errorf("invalid type %s", input.Type),
		)
		return
	}

	id, err := s.storage.Create(address, addressType, input.Description, ts)
	if err == storage.ErrExists {
		httputil.ResponseFailure(
			c,
			http.StatusConflict,
			err,
		)
		return
	} else if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s *Server) get(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/http/*Server.get")
	)

	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			fmt.Errorf("invalid id: %s", idVal),
		)
		return
	}

	logger = logger.With("id", id)
	logger.Debug("querying reserve address from database")

	ra, err := s.storage.Get(id)
	if err == storage.ErrNotExists {
		httputil.ResponseFailure(
			c,
			http.StatusNotFound,
			err,
		)
		return
	}
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, ra)
}
