package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type walletFeeQuery struct {
	httputil.TimeRangeQueryFreq
	ReserveAddr string `form:"reserve" binding:"required,isAddress"`
	WalletAddr  string `form:"walletAddr" binding:"required,isAddress"`
	Timezone    int64  `form:"timezone" binding:"isSupportedTimezone"`
}

func (ha *Server) getWalletFee(c *gin.Context) {
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

	_, _, err := query.Validate()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	fromTime := timeutil.TimestampMsToTime(query.From).UTC()
	toTime := timeutil.TimestampMsToTime(query.To).UTC()
	walletAddr := common.HexToAddress(query.WalletAddr).Hex()
	reserveAddr := common.HexToAddress(query.ReserveAddr).Hex()

	walletFee, err := ha.storage.GetAggregatedWalletFee(reserveAddr, walletAddr, query.Freq, fromTime, toTime, query.Timezone)

	if err != nil {
		ha.sugar.Errorw("reserve addr", query.ReserveAddr, "Wallet addr", query.WalletAddr,
			"from time", fromTime, "to time", toTime, "frequency", query.Freq)
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
		walletFee,
	)
}
