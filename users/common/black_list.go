package common

import (
	"encoding/json"
	"os"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	blackListFlag = "blacklist-file"
)

// NewBlacklistFlag returns flag for blacklist
func NewBlacklistFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   blackListFlag,
			Usage:  "json file to store blacklist",
			EnvVar: "BLACKLIST_FILE",
			Value:  "",
		},
	}
}

type Blacklist struct {
	bannedAddr map[ethereum.Address]struct{}
	logger     *zap.SugaredLogger
}

func newBlacklist(logger *zap.SugaredLogger, file string) (*Blacklist, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	var data []ethereum.Address
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return nil, errors.Wrap(err, "failed to decode json")
	}
	bannedAddr := make(map[ethereum.Address]struct{})
	for _, addr := range data {
		bannedAddr[addr] = struct{}{}
	}
	return &Blacklist{
		bannedAddr: bannedAddr,
		logger:     logger,
	}, nil
}

// NewBlacklistFromContext returns blacklist from cli.Context
func NewBlacklistFromContext(ctx *cli.Context, logger *zap.SugaredLogger) (*Blacklist, error) {
	return newBlacklist(logger, ctx.String(blackListFlag))
}

// IsBanned returns whether this address is banned
func (bl *Blacklist) IsBanned(address ethereum.Address) bool {
	_, ok := bl.bannedAddr[address]
	return ok
}
