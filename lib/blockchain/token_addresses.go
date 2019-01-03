package blockchain

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

type tokenAddresses struct {
	// ETHAddr is ethereum address
	ETHAddr common.Address
	// KNCAddr is KNC token address
	KNCAddr common.Address
	// WETHAddr is wrapped ETH address
	WETHAddr common.Address
	// KCCAddr is Kyber Community Coupon token address
	KCCAddr common.Address
}

var deploymentAddress = map[deployment.Deployment]tokenAddresses{
	deployment.Production: tokenAddresses{
		ETHAddr:  common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
		KNCAddr:  common.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200"),
		WETHAddr: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		KCCAddr:  common.HexToAddress("0x09677D0175DEC51E2215426Cddd055a71bf4228d"),
	},
	deployment.Staging: tokenAddresses{
		ETHAddr:  common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
		KNCAddr:  common.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200"),
		WETHAddr: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		KCCAddr:  common.HexToAddress("0x09677D0175DEC51E2215426Cddd055a71bf4228d"),
	},
	deployment.Ropsten: tokenAddresses{
		ETHAddr:  common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
		KNCAddr:  common.HexToAddress("0x4e470dc7321e84ca96fcaedd0c8abcebbaeb68c6"),
		WETHAddr: common.HexToAddress(""),
		KCCAddr:  common.HexToAddress(""),
	},
}

func getTokenAddressFromContext(c *cli.Context) (tokenAddresses, error) {
	deploymentMode, err := app.GetDeploymentModeFromContext(c)
	if err != nil {
		return tokenAddresses{}, err
	}
	tAddrs, ok := deploymentAddress[deploymentMode]
	if !ok {
		return tokenAddresses{}, fmt.Errorf("cannot get token addresses for deployment mode %s", deploymentMode.String())
	}
	return tAddrs, nil
}
