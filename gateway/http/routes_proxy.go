package http

//Option define initialize behaviour for server
type Option func(*Server) error

//WithTradeLogURL set TradeLogsProxy for server
func WithTradeLogURL(tradeLogsURL string) Option {
	return func(s *Server) error {
		tradeLogsProxyMW, err := newReverseProxyMW(tradeLogsURL)
		if err != nil {
			return err
		}
		s.r.GET("/trade-logs", tradeLogsProxyMW)
		s.r.GET("/trade-logs/:tx_hash", tradeLogsProxyMW)
		s.r.GET("/stats", tradeLogsProxyMW)
		s.r.GET("/top-tokens", tradeLogsProxyMW)
		s.r.GET("/top-integrations", tradeLogsProxyMW)
		s.r.GET("/top-reserves", tradeLogsProxyMW)
		s.r.GET("/big-trades", tradeLogsProxyMW)
		s.r.PUT("/big-trades", tradeLogsProxyMW)
		s.r.GET("/token-info", tradeLogsProxyMW)
		return nil
	}
}

//WithReserveRatesURL set resreve rate proxy for server
func WithReserveRatesURL(reserveRatesURL string) Option {
	return func(s *Server) error {
		reserveRateProxyMW, err := newReverseProxyMW(reserveRatesURL)
		if err != nil {
			return err
		}
		s.r.GET("/reserve-rates", reserveRateProxyMW)
		return nil
	}
}

//WithUserURL set user proxy for server
func WithUserURL(userURL string) Option {
	return func(s *Server) error {
		userProxyMW, err := newReverseProxyMW(userURL)
		if err != nil {
			return err
		}
		s.r.GET("/users", userProxyMW)
		s.r.POST("/users", userProxyMW)
		s.r.GET("/users-batch", userProxyMW)
		return nil
	}
}

//WithPriceAnalyticURL set price analytic proxy for server
func WithPriceAnalyticURL(priceAnalyticURL string) Option {
	return func(s *Server) error {
		priceProxyMW, err := newReverseProxyMW(priceAnalyticURL)
		if err != nil {
			return err
		}
		s.r.GET("/price-analytic-data", priceProxyMW)
		s.r.POST("/price-analytic-data", priceProxyMW)
		return nil
	}
}

//WithAppNamesURL set price analytic proxy for server
func WithAppNamesURL(appNamesURL string) Option {
	return func(s *Server) error {
		appNamesProxyMW, err := newReverseProxyMW(appNamesURL)
		if err != nil {
			return err
		}
		s.r.GET("/applications", appNamesProxyMW)
		s.r.POST("/applications", appNamesProxyMW)
		s.r.PUT("/applications", appNamesProxyMW)
		s.r.DELETE("/applications", appNamesProxyMW)
		return nil
	}
}

//WithCexTradesURL set cex trade proxy for server
func WithCexTradesURL(cexTradeURL string) Option {
	return func(s *Server) error {
		cexTradeURLMW, err := newReverseProxyMW(cexTradeURL)
		if err != nil {
			return err
		}
		s.r.GET("/trades", cexTradeURLMW)
		s.r.GET("/convert_to_eth_price", cexTradeURLMW)
		s.r.GET("/convert_trades", cexTradeURLMW)
		return nil
	}
}

//WithResreveAddressesURL set resreve addresses proxy for server
func WithResreveAddressesURL(reserveAddressesURL string) Option {
	return func(s *Server) error {
		reserveAddressURLMW, err := newReverseProxyMW(reserveAddressesURL)
		if err != nil {
			return err
		}
		s.r.POST("/addresses", reserveAddressURLMW)
		s.r.GET("/addresses/:id", reserveAddressURLMW)
		s.r.GET("/addresses", reserveAddressURLMW)
		s.r.PUT("/addresses/:id", reserveAddressURLMW)
		return nil
	}
}

// WithCexWithdrawalURL return withdraw proxy
func WithCexWithdrawalURL(cexWithdrawalURL string) Option {
	return func(s *Server) error {
		cexWithdrawalURLMW, err := newReverseProxyMW(cexWithdrawalURL)
		if err != nil {
			return err
		}
		s.r.GET("/withdrawals", cexWithdrawalURLMW)
		return nil
	}
}

// WithCexDepositURL return withdraw proxy
func WithCexDepositURL(cexDepositURL string) Option {
	return func(s *Server) error {
		cexDepositURLMW, err := newReverseProxyMW(cexDepositURL)
		if err != nil {
			return err
		}
		s.r.GET("/deposits", cexDepositURLMW)
		return nil
	}
}

//WithReserveTokenURL return reserve token proxy
func WithReserveTokenURL(reserveTokenURL string) Option {
	return func(s *Server) error {
		reserveTokenURLMW, err := newReverseProxyMW(reserveTokenURL)
		if err != nil {
			return err
		}
		s.r.GET("/reserve/tokens", reserveTokenURLMW)
		return nil
	}
}

//WithReserveTransactionURL return withdraw proxy
func WithReserveTransactionURL(reserveTransactionURL string) Option {
	return func(s *Server) error {
		reserveTransactionURLMW, err := newReverseProxyMW(reserveTransactionURL)
		if err != nil {
			return err
		}
		s.r.GET("/transactions", reserveTransactionURLMW)
		return nil
	}
}

//WithERC20APIURL return withdraw proxy
func WithERC20APIURL(erc20URL string) Option {
	return func(s *Server) error {
		erc20URLMW, err := newReverseProxyMW(erc20URL)
		if err != nil {
			return err
		}
		s.r.GET("/wallet/transactions", erc20URLMW)
		return nil
	}
}
