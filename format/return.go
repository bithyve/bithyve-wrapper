package format

import "sort"

// UtxoTxReturn is the return structure used in /utxotxs
type UtxoTxReturn struct {
	Utxos        [][]Utxo             `json:"Utxos"`
	Transactions []MultigetAddrReturn `json:"Txs"`
}

// Categorize does some nifty operations on the tx
func (tx *Tx) Categorize(InUseAddresses []string, ExternalAddresses []string) {
	var inputs = tx.Vin
	var outputs = tx.Vout
	var value, amountToSelf = float64(0), float64(0)
	var probableRecipientList []string
	var probableSenderList []string
	var selfRecipientList []string
	var selfSenderList []string

	for _, input := range inputs {
		var address = input.PrevOut.ScriptpubkeyAddress
		if len(address) == 0 {
			continue
		}
		if sort.SearchStrings(InUseAddresses, address) != 0 {
			value -= input.PrevOut.Value
			selfSenderList = append(selfSenderList, address)
		} else {
			probableSenderList = append(probableSenderList, address)
		}
	}

	for _, output := range outputs {
		var address = output.ScriptpubkeyAddress
		if len(address) == 0 {
			continue
		}
		if sort.SearchStrings(InUseAddresses, address) != 0 {
			value += output.Value
			if sort.SearchStrings(ExternalAddresses, address) != 0 {
				amountToSelf += output.Value
				selfRecipientList = append(selfRecipientList, address)
			}
		} else {
			probableRecipientList = append(probableRecipientList, address)
		}
	}

	if value > 0 {
		tx.TransactionType = "Received"
		tx.SenderAddresses = probableSenderList
	} else {
		if value+tx.Fee == 0 {
			tx.TransactionType = "Self"
			tx.SentAmount = amountToSelf + tx.Fee
			tx.ReceivedAmount = amountToSelf
			tx.SenderAddresses = selfSenderList
			tx.RecipientAddresses = selfRecipientList
		} else {
			tx.TransactionType = "Sent"
			tx.RecipientAddresses = probableRecipientList
		}
	}
}

// BalTxReturn is a struct used for the baltxs endpoint
type BalTxReturn struct {
	Balance      BalanceReturn
	Transactions []MultigetAddrReturn `json:"Txs"`
}

// TxReturn is used to return Txs
type TxReturn struct {
	Txs [][]Tx `json:"Txs"`
}

// BalanceReturn is a structure that is used for getting multiple balances
type BalanceReturn struct {
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
