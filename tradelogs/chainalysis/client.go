package chainalysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"go.uber.org/zap"

	ethereum "github.com/ethereum/go-ethereum/common"
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

func updateRegisterData(rd registerData, asset, txHash, receiverAdderss string) registerData {
	rd.RwData = append(rd.RwData, registerWithdrawal{
		Asset:   asset,
		Address: receiverAdderss,
	})
	rd.RstData = append(rd.RstData, registerSentTransfer{
		Asset:             asset,
		TransferReference: txHash,
	})
	return rd
}

// PushETHSentTransferEvent push eth sent transfer to chainalysis api
func (c *Client) PushETHSentTransferEvent(tradeLogs []common.TradeLog) error {
	mapRegisterData := make(map[ethereum.Address]registerData)
	for _, log := range tradeLogs {
		var (
			userAddress = log.UserAddress

			txHash          = log.TransactionHash.Hex()
			receiverAdderss = log.ReceiverAddress.Hex()
		)
		if strings.ToLower(log.DestAddress.Hex()) != ethAddress {
			continue
		}

		c.sugar.Debugw("sent transfer data",
			"user addr", userAddress,
			"receive addr", receiverAdderss,
			"tx hash", txHash)
		if rd, ok := mapRegisterData[userAddress]; ok {
			mapRegisterData[userAddress] = updateRegisterData(rd, ethSymbol, txHash, receiverAdderss)
		} else {
			mapRegisterData[userAddress] = registerData{
				RwData: []registerWithdrawal{
					{
						ethSymbol,
						receiverAdderss,
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
			c.sugar.Errorw("got error when register withdrawal address",
				"error", err.Error(),
				"user address", userAddress,
				"register withdrawal data", registerData.RwData)
			return err
		}
		if err := c.registerSentTransfer(userAddress, registerData.RstData); err != nil {
			c.sugar.Errorw("got error when register sent transfer",
				"error", err.Error(),
				"user address", userAddress,
				"register sent transfer data", registerData.RstData)
			return err
		}
	}
	return nil
}

// registerWithdrawalAddress register withdrawal address
func (c *Client) registerWithdrawalAddress(userAddr ethereum.Address, rw []registerWithdrawal) error {
	url := fmt.Sprintf("%s/users/%s/withdrawaladdresses", c.host, userAddr.Hex())
	c.sugar.Debugw("register withdrawal",
		"url", url,
		"body", rw)
	return c.registerChainAlysis(url, rw)
}

// registerSentTransfer register sent transfer
func (c *Client) registerSentTransfer(userAddr ethereum.Address, rst []registerSentTransfer) error {
	url := fmt.Sprintf("%s/users/%s/transfers/sent", c.host, userAddr.Hex())
	c.sugar.Debugw("register sent transfer",
		"url", url,
		"body", rst)
	return c.registerChainAlysis(url, rst)
}

// registerChainAlysis common function to register to chain alysis api
func (c *Client) registerChainAlysis(url string, data interface{}) error {
	body, err := json.Marshal(data)
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
	return nil
}
