package electrs

// GetBalanceFormat is a struct that us used to get the blanace from esplora
type GetBalanceFormat struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount float64 `json:"funded_txo_count"`
		FundedTxoSum   float64 `json:"funded_txo_sum"`
		SpentTxoCount  float64 `json:"spent_txo_count"`
		SpentTxoSum    float64 `json:"spent_txo_sum"`
		TxCount        float64 `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount float64 `json:"funded_txo_count"`
		FundedTxoSum   float64 `json:"funded_txo_sum"`
		SpentTxoCount  float64 `json:"spent_txo_count"`
		SpentTxoSum    float64 `json:"spent_txo_sum"`
		TxCount        float64 `json:"tx_count"`
	} `json:"mempool_stats"`
}
