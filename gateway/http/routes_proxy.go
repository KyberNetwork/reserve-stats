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
		s.r.GET("/burn-fee", tradeLogsProxyMW)
		s.r.GET("/asset-volume", tradeLogsProxyMW)
		s.r.GET("/reserve-volume", tradeLogsProxyMW)
		s.r.GET("/wallet-fee", tradeLogsProxyMW)
		s.r.GET("/user-volume", tradeLogsProxyMW)
		s.r.GET("/user-list", tradeLogsProxyMW)
		s.r.GET("/trade-summary", tradeLogsProxyMW)
		s.r.GET("/wallet-stats", tradeLogsProxyMW)
		s.r.GET("/country-stats", tradeLogsProxyMW)
		s.r.GET("/heat-map", tradeLogsProxyMW)
		s.r.GET("/integration-volume", tradeLogsProxyMW)
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

//WithCexTradesURL set cex trade proxy for server
func WithCexTradesURL(cexTradeURL string) Option {
	return func(s *Server) error {
		cexTradeURLMW, err := newReverseProxyMW(cexTradeURL)
		if err != nil {
			return err
		}
		s.r.GET("/cex_trades", cexTradeURLMW)
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

//WithCexWithdrawalURL return withdraw proxy
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