package http

import (
	"net/http"

	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type reserveVolumeQuery struct {
	From    uint64 `form:"from" `
	To      uint64 `form:"to"`
	Asset   string `form:"asset"`
	Freq    string `form:"freq"`
	Reserve string `form:"reserve" binding:"isAddress"`
}

func (sv *Server) getReserveVolume(c *gin.Context) {
	var (
		query  reserveVolumeQuery
		logger = sv.sugar.With("func", "tradelogs/volumehttp.getAssetVolume")
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	if !timeValidation(&query.From, &query.To, c, logger) {
		logger.Info("time validation returned invalid")
		return
	}
	token, err := sv.lookupToken(query.Asset)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	result, err := sv.storage.GetReserveVolume(ethereum.HexToAddress(query.Reserve), token, query.From, query.To, query.Freq)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}
