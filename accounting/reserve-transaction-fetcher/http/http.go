package http

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	txcommon "github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	maxTimeFrame     = time.Hour * 24 * 30 // 30 days
	defaultTimeFrame = time.Hour * 24      // 1 day
)

// Server is the HTTP server of accounting CEX getTrades HTTP API.
type Server struct {
	sugar *zap.SugaredLogger
	r     *gin.Engine
	host  string
	rts   storage.ReserveTransactionStorage
}

type getTransactionsQuery struct {
	httputil.TimeRangeQuery
	Types []string `form:"type"`
}

func (q *getTransactionsQuery) validate() (time.Time, time.Time, map[string]struct{}, error) {
	var types = make(map[string]struct{})

	fromTime, toTime, err := q.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}
	for _, typeString := range q.Types {
		_, ok := txcommon.TransactionTypes[typeString]
		if !ok {
			return time.Time{}, time.Time{}, nil, fmt.Errorf("invalid type %s", typeString)
		}
		types[typeString] = struct{}{}
	}
	//If the types is empty, return all types
	if len(types) == 0 {
		for typeString := range txcommon.TransactionTypes {
			types[typeString] = struct{}{}
		}

	}

	return fromTime, toTime, types, nil
}

type getTransactionsResponse struct {
	ERC20    []common.ERC20Transfer `json:"erc20,omitempty"`
	Normal   []common.NormalTx      `json:"normal,omitempty"`
	Internal []common.InternalTx    `json:"internal,omitempty"`
}

func (s *Server) getTransactions(c *gin.Context) {
	var (
		query       getTransactionsQuery
		erc20Txs    []common.ERC20Transfer
		normalTxs   []common.NormalTx
		internalTxs []common.InternalTx
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	from, to, types, err := query.validate()
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if _, ok := types[txcommon.ERC20.String()]; ok {
		erc20Txs, err = s.rts.GetERC20Transfer(from, to)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}
	}
	if _, ok := types[txcommon.Normal.String()]; ok {
		normalTxs, err = s.rts.GetNormalTx(from, to)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}
	}
	if _, ok := types[txcommon.Internal.String()]; ok {
		internalTxs, err = s.rts.GetInternalTx(from, to)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}
	}
	c.JSON(http.StatusOK, getTransactionsResponse{
		ERC20:    erc20Txs,
		Normal:   normalTxs,
		Internal: internalTxs,
	})
}

// NewServer creates a new instance of Server.
func NewServer(logger *zap.Logger, host string, rts storage.ReserveTransactionStorage) *Server {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	sugar := logger.Sugar()
	return &Server{sugar: sugar, r: r, host: host, rts: rts}
}

func (s *Server) register() {
	s.r.GET("/transactions", s.getTransactions)
}

// Run starts the HTTP server and runs in foreground until terminate by user.
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
