package http

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/gin-gonic/gin"
)

type assetVolumeQuery struct {
	From  uint64 `form:"from" `
	To    uint64 `form:"to"`
	Asset string `form:"asset"`
	Freq  string `form:"freq"`
}

//Server struct to represent a http server service
type Server struct {
	sugar   *zap.SugaredLogger
	r       *gin.Engine
	host    string
	storage storage.VolumeStorage
	setting coreSetting
}

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, storage storage.VolumeStorage, host string, sett coreSetting) *Server {
	r := gin.Default()
	return &Server{
		sugar:   sugar,
		storage: storage,
		r:       r,
		host:    host,
		setting: sett,
	}
}

func (sv *Server) getAssetVolume(c *gin.Context) {
	var (
		query  assetVolumeQuery
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

func (sv *Server) register() {
	sv.r.GET("/asset-volume", sv.getAssetVolume)
}

//Run start server and serve
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
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

func (sv *Server) lookupToken(ID string) (core.Token, error) {
	tokens, err := sv.setting.GetActiveTokens()
	if err != nil {
		return core.Token{}, err
	}
	for _, token := range tokens {
		if token.ID == ID {
			return token, nil
		}
	}
	return core.Token{}, fmt.Errorf("cannot find token %s in current core setting", ID)
}
