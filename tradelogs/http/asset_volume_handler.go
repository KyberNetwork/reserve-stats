package http

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/gin-gonic/gin"
)

type assetVolumeQuery struct {
	From  uint64 `form:"from" `
	To    uint64 `form:"to"`
	Asset string `form:"asset"`
	Freq  string `form:"freq"`
}

func (sv *Server) getAssetVolume(c *gin.Context) {
	var (
		query  assetVolumeQuery
		logger = sv.sugar.With("func", "tradelogs/http/Server.getAssetVolume")
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
	token, err := core.LookupToken(sv.coreSetting, query.Asset)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	result, err := sv.storage.GetAssetVolume(token, query.From, query.To, query.Freq)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

func timeValidation(fromTime, toTime *uint64, c *gin.Context, logger *zap.SugaredLogger) bool {
	if *fromTime == 0 && *toTime == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid time frame query, from: %d, to: %d", *fromTime, *toTime)},
		)
		return false
	}
	now := time.Now().UTC()
	if *toTime == 0 {
		*toTime = timeutil.TimeToTimestampMs(now)
		logger.Debug("using default to query time", "to", *toTime)

		if *fromTime == 0 {
			*fromTime = timeutil.TimeToTimestampMs(now.Add(-time.Hour))
			logger = logger.With("from", *fromTime)
			logger.Debug("using default from query time", "from", *fromTime)
		}
	}
	return true
}
