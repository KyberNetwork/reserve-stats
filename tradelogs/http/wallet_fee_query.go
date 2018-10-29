package http

import (
	"net/http"

	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/gin-gonic/gin"
)

type walletFeeQuery struct {
	From        uint64 `form:"from" binding:"required"`
	To          uint64 `form:"to" binding:"required"`
	Freq        string `form:"freq" binding:"required,isFreq"`
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

	fromTime := timeutil.TimestampMsToTime(query.From)
	toTime := timeutil.TimestampMsToTime(query.To)

	walletFee, err := ha.storage.GetAggregatedWalletFee(query.ReserveAddr, query.WalletAddr, query.Freq, fromTime, toTime, query.Timezone)

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
