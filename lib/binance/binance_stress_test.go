package binance

import (
	"os"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestBinance(t *testing.T) {
	t.Skip("test disabled as require binance API key pair")
	var (
		apiKey    = os.Getenv("BINANCE_API_KEY")
		secretKey = os.Getenv("BINANCE_SECRET_KEY")
	)

	sugar := testutil.MustNewDevelopmentSugaredLogger()

	bn := NewBinance(apiKey, secretKey, sugar)

	const requests = 10000
	var g errgroup.Group
	for i := 0; i < requests; i++ {
		i := i
		g.Go(func() error {
			i := i
			if _, err := bn.GetExchangeInfo(); err != nil {
				return err
			}
			sugar.Infow("request completed", "order", i)
			return nil
		})
	}

	assert.NoError(t, g.Wait())
}
