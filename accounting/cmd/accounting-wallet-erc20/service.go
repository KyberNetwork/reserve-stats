package main

import (
	"database/sql"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/wallet-erc20/storage"
)

type erc20TransferFetcher interface {
	ERC20Transfer(addr ethereum.Address, from, to *big.Int) ([]common.ERC20Transfer, error)
}

type service struct {
	sugar *zap.SugaredLogger
	f     erc20TransferFetcher
	db    storage.Interface

	addrs     []ethereum.Address
	fromBlock *big.Int
	toBlock   *big.Int
}

func newService(
	sugar *zap.SugaredLogger,
	f erc20TransferFetcher,
	db storage.Interface,
	addrs []ethereum.Address,
	fromBlock, toBlock *big.Int) *service {
	return &service{sugar: sugar, f: f, db: db, addrs: addrs, fromBlock: fromBlock, toBlock: toBlock}
}

func (s *service) run() error {
	var logger = s.sugar.With("func", "service.run")

	for _, walletAddr := range s.addrs {
		var (
			fromBlock *big.Int
			toBlock   *big.Int
		)
		if s.fromBlock != nil {
			fromBlock = big.NewInt(0).Set(s.fromBlock)
		}
		if s.toBlock != nil {
			toBlock = big.NewInt(0).Set(s.toBlock)
		}

		if fromBlock == nil {
			lastStoredBlock, err := s.db.GetLastStoredBlock(walletAddr)
			switch err {
			case sql.ErrNoRows:
				logger.Infow("no record found, fetching from beginning",
					"wallet", walletAddr)
			case nil:
				logger.Infow("fetching from last block", "last_block", lastStoredBlock+1)
				fromBlock = big.NewInt(int64(lastStoredBlock + 1))
			default:
				return err
			}
		}
		transfers, err := s.f.ERC20Transfer(walletAddr, fromBlock, toBlock)
		if err != nil {
			return err
		}
		if err := s.db.UpdateERC20Transfers(transfers); err != nil {
			return err
		}
	}

	return nil
}
