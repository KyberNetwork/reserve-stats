package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type walletStatsQuery struct {
	httputil.TimeRangeQuery
	WalletAddr string `form:"walletAddr,isEthereumAddress"`
	Timezone   int8   `form:"timezone" binding:"isSupportTimezone"`
}

func (sv *Server) getWalletStats(c *gin.Context) {
	var (
		query walletStatsQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	walletAddr := ethereum.HexToAddress(query.WalletAddr)
	from := timeutil.TimestampMsToTime(query.From)
	to := timeutil.TimestampMsToTime(query.To)
	walletStats, err := sv.storage.GetWalletStats(from, to, walletAddr.Hex(), query.Timezone)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, walletStats)
}
