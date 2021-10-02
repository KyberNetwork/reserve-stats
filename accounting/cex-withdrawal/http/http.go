package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	huobiStorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	maxTimeFrame     = time.Hour * 24 * 30 // 30 days
	defaultTimeFrame = time.Hour * 24      // 1 day
)

// Server is the engine to serve cex-trade-withdrawal API query
type Server struct {
	r         *gin.Engine
	huobiDB   huobiStorage.Interface
	binanceDB withdrawalstorage.Interface
	host      string
	sugar     *zap.SugaredLogger
}

type queryInput struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex" binding:"dive,isValidCEXName"`
}

// BinanceWithdrawalResponse ...
type BinanceWithdrawalResponse struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Address   string  `json:"address"`
	Asset     string  `json:"asset"`
	TxID      string  `json:"txId"`
	ApplyTime int64   `json:"applyTime"`
	Status    int64   `json:"status"`
	TxFee     float64 `json:"transactionFee"`
}

type response struct {
	Huobi   map[string][]huobi.WithdrawHistory     `json:"huobi,omitempty"`
	Binance map[string][]BinanceWithdrawalResponse `json:"binance,omitempty"`
}

func (sv *Server) get(c *gin.Context) {
	var (
		query            queryInput
		logger           = sv.sugar.With("func", caller.GetCurrentFunctionName())
		huobiWithdrawals = make(map[string][]huobi.WithdrawHistory)
		// binanceWithdrawals []binance.WithdrawHistory
		binanceResponse = make(map[string][]BinanceWithdrawalResponse)
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if len(query.Exchanges) == 0 {
		query.Exchanges = []string{
			common.Huobi.String(),
			common.Binance.String()}
	}

	from, to, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	logger = logger.With("from", from, "to", to, "exchanges", query.Exchanges)
	logger.Debug("querying withdrawals from database")

	for _, cex := range query.Exchanges {
		switch cex {
		case common.Huobi.String():
			huobiWithdrawals, err = sv.huobiDB.GetWithdrawHistory(from, to)
			if err != nil {
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
		case common.Binance.String():
			binanceWithdrawals, err := sv.binanceDB.GetWithdrawHistory(from, to)
			if err != nil {
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			for account, withdrawals := range binanceWithdrawals {
				response := []BinanceWithdrawalResponse{}
				for _, withdrawal := range withdrawals {
					applyTime, sErr := time.Parse("2006-01-02 15:04:05", withdrawal.ApplyTime)
					if sErr != nil {
						return
					}
					amount, err := strconv.ParseFloat(withdrawal.Amount, 64)
					if err != nil {
						sv.sugar.Errorw("failed to parse withdraw amount", "error", err)
						httputil.ResponseFailure(
							c,
							http.StatusInternalServerError,
							err,
						)
						return
					}
					txFee, err := strconv.ParseFloat(withdrawal.TxFee, 64)
					if err != nil {
						sv.sugar.Errorw("failed to parse tx fee", "error", err)
					}
					response = append(response, BinanceWithdrawalResponse{
						ID:        withdrawal.ID,
						Amount:    amount,
						Address:   withdrawal.Address,
						Asset:     withdrawal.Asset,
						TxID:      withdrawal.TxID,
						ApplyTime: applyTime.UnixMilli(),
						Status:    withdrawal.Status,
						TxFee:     txFee,
					})
				}
				binanceResponse[account] = response
			}
		}
	}

	c.JSON(http.StatusOK, response{
		Huobi:   huobiWithdrawals,
		Binance: binanceResponse,
	})
}

func (sv *Server) register() {
	sv.r.GET("/withdrawals", sv.get)
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.register()
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, huobiDB huobiStorage.Interface, binanceDB withdrawalstorage.Interface, sugar *zap.SugaredLogger) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:         r,
		huobiDB:   huobiDB,
		binanceDB: binanceDB,
		host:      host,
		sugar:     sugar,
	}, nil
}
