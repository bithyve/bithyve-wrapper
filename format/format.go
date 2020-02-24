package format

// RequestFormat is the format in which incoming requests hsould arrive for the wrapper to process
type RequestFormat struct {
	Addresses []string `json:"addresses"`
}

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

// TxReturn is used to return Txs
type TxReturn struct {
	Txs [][]Tx `json:"Txs"`
}

// MultigetBalanceReturn is a structure that is used for getting multiple balances
type MultigetBalanceReturn struct {
	Balance            float64
	UnconfirmedBalance float64
}

// MultigetAddrReturn is a structure used for multiple addresses json return
type MultigetAddrReturn struct {
	TotalTransactions       float64
	ConfirmedTransactions   float64
	UnconfirmedTransactions float64
	Transactions            []Tx
	Address                 string
}

// UtxoVout is a structure for output utxos
type UtxoVout struct {
	Scriptpubkey        string  `json:"scriptpubkey"`
	ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
	ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
	ScriptpubkeyType    string  `json:"scriptpubkey_type"`
	Value               float64 `json:"value"`
	Index               int
	Address             string
}

// Utxo is a copy of the esplora utxo struct
type Utxo struct {
	Txid   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed   bool    `json:"confirmed"`
		BlockHeight float64 `json:"block_height"`
		BlockHash   string  `json:"block_hash"`
		BlockTime   float64 `json:"block_time"`
	} `json:"status"`
	Value   float64 `json:"value"`
	Address string
}

// Tx is a copy of the transaction structure used by esplora
type Tx struct {
	Txid     string  `json:"txid"`
	Version  float64 `json:"version"`
	Locktime float64 `json:"locktime"`
	Vin      []struct {
		Txid    string  `json:"txid"`
		Vout    float64 `json:"vout"`
		PrevOut struct {
			Scriptpubkey        string  `json:"scriptpubkey"`
			ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
			ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
			ScriptpubkeyType    string  `json:"scriptpubkey_type"`
			Value               float64 `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     float64  `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string  `json:"scriptpubkey"`
		ScriptpubkeyAsm     string  `json:"scriptpubkey_asm"`
		ScriptpubkeyAddress string  `json:"scriptpubkey_address"`
		ScriptpubkeyType    string  `json:"scriptpubkey_type"`
		Value               float64 `json:"value"`
	}
	Size   float64 `json:"size"`
	Weight float64 `json:"weight"`
	Fee    float64 `json:"fee"`
	Status struct {
		Confirmed   bool    `json:"confirmed"`
		BlockHeight float64 `json:"block_height"`
		BlockHash   string  `json:"block_hash"`
		BlockTime   float64 `json:"block_time"`
	}
	NumberofConfirmations float64
}

// FeeResponse is a struct that is returned when a fee query is made
type FeeResponse struct {
	Two              float64 `json:"2"`
	Three            float64 `json:"3"`
	Four             float64 `json:"4"`
	Six              float64 `json:"6"`
	Ten              float64 `json:"10"`
	Twenty           float64 `json:"20"`
	OneFourFour      float64 `json:"144"`
	FiveZeroFour     float64 `json:"504"`
	OneThousandEight float64 `json:"1008"`
}

// BalTx is a struct used for the baltxs endpoint
type BalTx struct {
	Balance      MultigetBalanceReturn
	Transactions []MultigetAddrReturn `json:"Txs"`
}
