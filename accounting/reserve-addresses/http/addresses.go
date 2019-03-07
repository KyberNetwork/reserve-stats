package http

import (
	"fmt"
	"net/http"
	"strconv"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
)

type createInput struct {
	Address     string `json:"address" binding:"required,isAddress"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) create(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/http/*Server.create")
		input  createInput
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

	addressType, ok := common.IsValidAddressType(input.Type)
	if !ok {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			fmt.Errorf("invalid type %s", input.Type),
		)
		return
	}

	id, err := s.storage.Create(address, addressType, input.Description)
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

// getIDParam gets and validates the id parameter from given context.
func getIDParam(c *gin.Context) (uint64, error) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idVal)
	}
	return id, nil
}

func (s *Server) get(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "accounting/reserve-addresses/http/*Server.get")
	)

	id, err := getIDParam(c)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
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

func (s *Server) getAll(c *gin.Context) {
	addrs, err := s.storage.GetAll()
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	if addrs == nil {
		addrs = []*common.ReserveAddress{}
	}

	c.JSON(http.StatusOK, addrs)
}

type updateInput struct {
	Address     string  `json:"address" binding:"isAddress"`
	Type        *string `json:"type"`
	Description string  `json:"description"`
}

func (s *Server) update(c *gin.Context) {
	var (
		logger      = s.sugar.With("func", "accounting/reserve-addresses/http/*Server.update")
		input       updateInput
		addressType *common.AddressType
	)

	id, err := getIDParam(c)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
	}

	logger = logger.With("id", id)

	if err = c.ShouldBindJSON(&input); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if input.Type != nil {
		validType, ok := common.IsValidAddressType(*input.Type)
		if !ok {
			httputil.ResponseFailure(
				c,
				http.StatusBadRequest,
				fmt.Errorf("invalid type %s", *input.Type),
			)
			return
		}
		addressType = &validType
	}

	logger = logger.With(
		"address", input.Address,
		"type", input.Type,
		"description", input.Description,
	)

	logger.Debug("updating reserve address")

	address := ethereum.HexToAddress(input.Address)

	if err = s.storage.Update(id, address, addressType, input.Description); err == storage.ErrNotExists {
		httputil.ResponseFailure(
			c,
			http.StatusNotFound,
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
	c.Status(http.StatusNoContent)
}
