package main

type duneQueryData struct {
	QueryResult struct {
		Data struct {
			Rows []struct {
				CallTxHash  string `json:"call_tx_hash"`
				WalletID    string `json:"walletId"`
				Src         string `json:"src"`
				Dest        string `json:"dest"`
				DestAddress string `json:"destAddress"`
				CallSuccess bool   `json:"call_success"`
			} `json:"rows"`
		} `json:"data"`
	} `json:"query_result"`
}
