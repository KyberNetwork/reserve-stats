package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type walletFeeQuery struct {
	httputil.TimeRangeQueryFreq
	ReserveAddr string `form:"reserve" binding:"required,isAddress"`
	WalletAddr  string `form:"walletAddr" binding:"required,isAddress"`
	Timezone    int8   `form:"timezone" binding:"isSupportedTimezone"`
}

func (sv *Server) getWalletFee(c *gin.Context) {
	var query walletFeeQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	fromTime, toTime, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	// normalize Ethereum addresses
	walletAddr := common.HexToAddress(query.WalletAddr).Hex()
	reserveAddr := common.HexToAddress(query.ReserveAddr).Hex()

	walletFee, err := sv.storage.GetAggregatedWalletFee(reserveAddr, walletAddr, query.Freq, fromTime, toTime, query.Timezone)
	if err != nil {
		sv.sugar.Errorw("reserve addr", query.ReserveAddr, "Wallet addr", query.WalletAddr,
			"from time", fromTime, "to time", toTime, "frequency", query.Freq)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		walletFee,
	)
}
