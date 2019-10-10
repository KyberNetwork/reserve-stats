package chainalysis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	timeout = time.Minute * 5

	ethSymbol  = "ETH"
	ethAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

// Client is implementation of chainalysis client
type Client struct {
	host   string
	apiKey string
	sugar  *zap.SugaredLogger
	client *http.Client
}

// NewChainAlysisClient creates a new chainalysis client instance.
func NewChainAlysisClient(sugar *zap.SugaredLogger, host, apiKey string) *Client {
	c := &Client{
		host:   host,
		apiKey: apiKey,
		sugar:  sugar,
		client: &http.Client{Timeout: timeout},
	}
	return c
}

type registerData struct {
	RwData  []registerWithdrawal
	RstData []registerSentTransfer
}

type registerWithdrawal struct {
	Asset   string `json:"asset"`
	Address string `json:"address"`
}

type registerSentTransfer struct {
	Asset             string `json:"asset"`
	TransferReference string `json:"transferReference"`
}

type errorResponse struct {
	Message string
}

func updateRegisterData(rd registerData, asset, txHash, receiverAddress string) registerData {
	rd.RwData = append(rd.RwData, registerWithdrawal{
		Asset:   asset,
		Address: receiverAddress,
	})
	rd.RstData = append(rd.RstData, registerSentTransfer{
		Asset:             asset,
		TransferReference: txHash,
	})
	return rd
}

// PushETHSentTransferEvent push eth sent transfer to chainalysis api
func (c *Client) PushETHSentTransferEvent(tradeLogs []common.TradeLog) error {
	var (
		logger = c.sugar.With("func", caller.GetCurrentFunctionName())

		mapRegisterData = make(map[ethereum.Address]registerData)
	)
	for _, log := range tradeLogs {
		var (
			userAddress = log.UserAddress

			txHash          = log.TransactionHash.Hex()
			receiverAddress = log.ReceiverAddress.Hex()
		)
		if strings.ToLower(log.DestAddress.Hex()) != ethAddress {
			continue
		}

		if rd, ok := mapRegisterData[userAddress]; ok {
			mapRegisterData[userAddress] = updateRegisterData(rd, ethSymbol, txHash, receiverAddress)
		} else {
			mapRegisterData[userAddress] = registerData{
				RwData: []registerWithdrawal{
					{
						ethSymbol,
						receiverAddress,
					},
				},
				RstData: []registerSentTransfer{
					{
						ethSymbol,
						txHash,
					},
				},
			}
		}
	}
	for userAddress, registerData := range mapRegisterData {
		if err := c.registerWithdrawalAddress(userAddress, registerData.RwData); err != nil {
			logger.Errorw("got error when register withdrawal address", "error", err.Error())
			return err
		}
		if err := c.registerSentTransfer(userAddress, registerData.RstData); err != nil {
			logger.Errorw("got error when register sent transfer", "error", err.Error())
			return err
		}
	}
	return nil
}

// registerWithdrawalAddress register withdrawal address
func (c *Client) registerWithdrawalAddress(userAddr ethereum.Address, rw []registerWithdrawal) error {
	var (
		logger = c.sugar.With("func", caller.GetCurrentFunctionName())

		url = fmt.Sprintf("%s/users/%s/withdrawaladdresses", c.host, userAddr.Hex())
	)
	logger.Debugw("register withdrawal",
		"url", url,
		"user address", userAddr,
		"body", rw)
	return c.registerChainAlysis(url, rw)
}

// registerSentTransfer register sent transfer
func (c *Client) registerSentTransfer(userAddr ethereum.Address, rst []registerSentTransfer) error {
	var (
		logger = c.sugar.With("func", caller.GetCurrentFunctionName())
		url    = fmt.Sprintf("%s/users/%s/transfers/sent", c.host, userAddr.Hex())
	)
	logger.Debugw("register sent transfer",
		"url", url,
		"user address", userAddr,
		"body", rst)
	return c.registerChainAlysis(url, rst)
}

// registerChainAlysis common function to register to chain alysis api
func (c *Client) registerChainAlysis(url string, data interface{}) error {
	var logger = c.sugar.With("func", caller.GetCurrentFunctionName())
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Token", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()
	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info("register successfully")
	case http.StatusBadRequest, http.StatusForbidden:
		eResp := errorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&eResp); err != nil {
			return err
		}
		return errors.New(eResp.Message)
	default:
		return fmt.Errorf("got unexpected code: %d", resp.StatusCode)
	}
	return nil
}
