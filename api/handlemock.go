package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/btcsuite/btcd/wire"
)

type TransactionStore struct {
	ScriptHash string      `json:"scripthash"`
	TxDetails  interface{} `json:"transaction"`
}

func handleGetMockHeader(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	datafilePath := filepath.Join("mock_tx", "header.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read Balance data: %v", err))
	}
	var txStre []TransactionStore
	var balance map[string]interface{}
	err = json.Unmarshal(data, &txStre)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse TransactionStore of Balance data: %v", err))
	}
	for _, txs := range txStre {
		if scripthash == txs.ScriptHash {
			balData, err := json.Marshal(txs.TxDetails)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Marshel Balance data: %v", err))
			}

			err = json.Unmarshal(balData, &balance)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		} else {
			balData := "{\"confirmed\": 0,\"unconfirmed\": 0}"

			err = json.Unmarshal([]byte(balData), &balance)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		}
	}

	return makeResponse(id, balance, nil)
}
func handleGetMockBalance(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	datafilePath := filepath.Join("mock_tx", "balance.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read Balance data: %v", err))
	}
	var txStre []TransactionStore
	var balance map[string]interface{}
	err = json.Unmarshal(data, &txStre)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse TransactionStore of Balance data: %v", err))
	}
	for _, txs := range txStre {
		if scripthash == txs.ScriptHash {
			balData, err := json.Marshal(txs.TxDetails)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Marshel Balance data: %v", err))
			}

			err = json.Unmarshal(balData, &balance)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		} else {
			balData := "{\"confirmed\": 0,\"unconfirmed\": 0}"

			err = json.Unmarshal([]byte(balData), &balance)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		}
	}

	return makeResponse(id, balance, nil)
}

func handleMockGetHistory(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	datafilePath := filepath.Join("mock_tx", "txshistory.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read Tx History data: %v", err))
	}
	var txStre []TransactionStore
	var txhistory []map[string]interface{}
	err = json.Unmarshal(data, &txStre)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Tx History data: %v", err))
	}
	for _, txs := range txStre {
		if scripthash == txs.ScriptHash {
			balData, err := json.Marshal(txs.TxDetails)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Marshel Balance data: %v", err))
			}

			err = json.Unmarshal(balData, &txhistory)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		} else {
			balData := "[]"

			err = json.Unmarshal([]byte(balData), &txhistory)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		}
	}
	// For demonstration, returning a mock response

	return makeResponse(id, txhistory, nil)
}
func handleMockListUnspent(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	datafilePath := filepath.Join("mock_tx", "utxo.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read UTXO data: %v", err))
	}
	var txStre []TransactionStore
	var utxos []map[string]interface{}
	err = json.Unmarshal(data, &txStre)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse UTXO data: %v", err))
	}
	for _, txs := range txStre {
		if scripthash == txs.ScriptHash {
			balData, err := json.Marshal(txs.TxDetails)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Marshel Balance data: %v", err))
			}

			err = json.Unmarshal(balData, &utxos)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		} else {
			balData := "[]"

			err = json.Unmarshal([]byte(balData), &utxos)
			if err != nil {
				return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
			}
		}
	}

	return makeResponse(id, utxos, nil)
}
func handleMockMempoolFee(id interface{}) string {
	datafilePath := filepath.Join("mock_tx", "feehistogram.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read Fee Histogram data: %v", err))
	}
	var feeHist [][]interface{}
	err = json.Unmarshal(data, &feeHist)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Fee Histogram data: %v", err))
	}
	return makeResponse(id, feeHist, nil)
}
func handleMockHandle(id interface{}) string {
	datafilePath := filepath.Join("mock_tx", "feehistogram.json")
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to read Fee Histogram data: %v", err))
	}
	var feeHist [][]interface{}
	err = json.Unmarshal(data, &feeHist)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Fee Histogram data: %v", err))
	}
	return makeResponse(id, feeHist, nil)
}
func estimatedFee(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}

	//_, ok := params[0].(int)

	return makeResponse(id, 0.00001, nil)
}
func broadcastTransaction(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	rawtx, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid TransactionHash")
	}
	txBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		return "" //, fmt.Errorf("failed to decode raw transaction: %v", err)
	}

	tx := wire.NewMsgTx(2)
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return "" //, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	txID := tx.TxHash().String()
	return makeResponse(id, txID, nil)
}
