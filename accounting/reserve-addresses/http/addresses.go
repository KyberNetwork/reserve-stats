package http

import (
	"fmt"
	"net/http"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
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
	if err != nil {
		// check if return error is a known pq error
		pErr, ok := err.(*pq.Error)
		if !ok {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}

		// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
		// 23505: unique_violation
		if pErr.Code == "23505" {
			httputil.ResponseFailure(
				c,
				http.StatusConflict,
				pErr,
			)
			return
		}

		// unknown pq error
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			pErr,
		)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
