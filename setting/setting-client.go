package setting

import "github.com/KyberNetwork/reserve-stats/common"

type SettingClient struct {
	coreURL string
}

func NewSettingClient(coreURL string) (*SettingClient, error) {
	return &SettingClient{}, nil
}

func (sc *SettingClient) GetActiveTokens() ([]common.Token, error) {
	result := []common.Token{}
	return result, nil
}

func (sc *SettingClient) GetInternalTokens() ([]common.Token, error) {
	result := []common.Token{}
	return result, nil
}
