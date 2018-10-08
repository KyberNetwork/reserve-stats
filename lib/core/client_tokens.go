package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type allSettingsResponse struct {
	commonResponse
	Data tokensData `json:"data"`
}

type tokensData struct {
	Tokens *tokenList `json:"tokens"`
}

type tokenList struct {
	Tokens []Token `json:"tokens"`
}

// generateNonce returns nonce header required to use Core API,
// which is current timestamp in milliseconds.
func generateNonce() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(now, 10)
}

// Tokens returns all configured tokens.
// Example response JSON:
// {
//  "data": {
//    "tokens": {
//      "tokens": [
//        {
//          "id": "ABT",
//          "name": "ArcBlock",
//          "address": "0xb98d4c97425d9908e66e53a6fdf673acca0be986",
//          "decimals": 18,
//          "active": true,
//          "internal": true,
//          "last_activation_change": 1535021910190
//        },
//        {
//          "id": "ADX",
//          "name": "AdEx",
//          "address": "0x4470BB87d77b963A013DB939BE332f927f2b992e",
//          "decimals": 4,
//          "active": true,
//          "internal": false,
//          "last_activation_change": 1535021910195
//        }
//      ]
//    }
//  }
//}
func (c *Client) Tokens() ([]Token, error) {
	const endpoint = "/setting/all-settings"
	var params = make(map[string]string)

	params["nonce"] = generateNonce()

	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}

	var settingsResponse = &allSettingsResponse{}
	if err = json.NewDecoder(rsp.Body).Decode(&settingsResponse); err != nil {
		return nil, err
	}

	if settingsResponse.Success != true {
		return nil, fmt.Errorf("got an error from server: %s", settingsResponse.Reason)
	}

	return settingsResponse.Data.Tokens.Tokens, nil
}

func (c *Client) getTokens(endpoint string) ([]Token, error) {
	var params = make(map[string]string)
	params["nonce"] = generateNonce()
	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	var tokenReply = &TokensReply{}
	if err = json.NewDecoder(rsp.Body).Decode(&tokenReply); err != nil {
		return nil, err
	}
	return tokenReply.Data, nil
}

// GetInternalTokens return list of internal token from Kyber reserve
func (c *Client) GetInternalTokens() ([]Token, error) {
	const endpoint = "setting/internal-tokens"
	return c.getTokens(endpoint)
}

// GetActiveTokens return list of active token from external reserve
func (c *Client) GetActiveTokens() ([]Token, error) {
	const endpoint = "setting/active-tokens"
	return c.getTokens(endpoint)
}
