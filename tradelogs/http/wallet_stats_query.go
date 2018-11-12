package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type walletStatsQuery struct {
	httputil.TimeRangeQuery
	WalletAddr string `form:"walletAddr,isEthereumAddress"`
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
	walletStats, err := sv.storage.GetWalletStats(query.From, query.To, walletAddr.Hex())
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
