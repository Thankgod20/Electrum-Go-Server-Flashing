package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	//"github.com/btcsuite/btcd/wire"
)

type TransactionStore struct {
	ScriptHash string      `json:"scripthash"`
	TxDetails  interface{} `json:"transaction"`
}

func (sc *Server) handleGetHeader(id interface{}, params []interface{}) string {
	hex, height := sc.bitcoinAPI.GetBlockHeader()

	var balance map[string]interface{}

	balData := "{\"height\":" + height + ",\"hex\":\"" + hex + "\"}"

	err := json.Unmarshal([]byte(balData), &balance)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
	}

	return makeResponse(id, balance, nil)
}
func (sc *Server) handleGetBalance(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	bal := sc.bitcoinAPI.GetBalance(scripthash)
	//log.Println("Balances", bal)
	var balance map[string]interface{}
	err := json.Unmarshal([]byte(bal), &balance)
	if err != nil {
		return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
	}
	return makeResponse(id, balance, nil)
}
func (sc *Server) handleGetTransaction(id interface{}, params []interface{}) string {
	if len(params) != 2 {
		return makeResponse(id, nil, "invalid params")
	}
	trxID, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	if params[1].(bool) {
		_, rawtrx := sc.bitcoinAPI.GetTransaction(trxID)
		return makeResponse(id, rawtrx, nil)
	} else {
		trxhex, _ := sc.bitcoinAPI.GetTransaction(trxID)

		return makeResponse(id, trxhex, nil)
	}
}
func (sc *Server) handleTrxGetHistory(id interface{}, params []interface{}) string {

	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	var transactions []map[string]interface{}
	trans, _ := sc.bitcoinAPI.GetTransactionHistory(scripthash)
	//fmt.Println("History", trans, mempool)
	//	txhisty := sc.bitcoinAPI.GetTransactionHistory(scripthash)
	//	var txhistory []map[string]interface{}

	// For demonstration, returning a mock response
	if len(trans) > 0 {
		transactions = trans
	} else {
		balData := "[]"

		err := json.Unmarshal([]byte(balData), &transactions)
		if err != nil {
			return makeResponse(id, nil, fmt.Sprintf("Failed to Parse Balance data: %v", err))
		}
	}
	return makeResponse(id, transactions, nil)
}
func (sc *Server) handleListUnspent(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	utxos := sc.bitcoinAPI.GetTransactionUTXOs(scripthash)
	log.Println("UTXOS", utxos)

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
	block, ok := params[0].(float64)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	//fmt.Println((params[0].(float64)))
	//_, ok := params[0].(int)
	fees := (0.0001 * block) //float64(scripthash))
	return makeResponse(id, fees, nil)
}

func (sc *Server) broadcastTransaction(id interface{}, params []interface{}) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	rawtx, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid TransactionHash")
	}
	txID := sc.bitcoinAPI.BroadCastTX(rawtx)
	return makeResponse(id, txID, nil)
}

//02000000000101572ce1c3d8818d445883e0372a0b5923b1f017a4056dfeceeaf92a702a1aae620000000000000000800280841e0000000000160014b863ae8777f8387cfeb2f4424503616ab4a7841b827f902f00000000160014f26e60250c5d52753a46e86796d24870c9b51c4b02483045022100a1edc977b15680549f4ac69cc596e0bc32a88c750af4417153558358b6d77c81022017250d0672a543c42abbb955638baf8d06e759cc0525b44f704db7649b568cb40121031b08427f3a17132260d47a9b0b2edc2b31092c3ddff639e6c1121ed6e952e1e500000000

//02000000000101c01e1d8338c72a842a6ae890dbb48543982c2a5ddbccaf3783b7b87529411ec10000000000000000800200e1f505000000001600145ca41ac1cb90764cd84e4d7dec0ad510daaf07f702e5a4350000000016001445c35beb1f7d47ef8fa2501838b0e8f3b0b5a8fd02473044022068e9667e43d4b6ac4fa1b438047e4d296bf1ca1726cf95f3a09cfd4b8362de2e0220221419f8fdef59f55c843b9e84c02ff3d6b7643c25aa9fc417d76e311e60fd28012103d664dfb8818daaee9e21d1acdecb70801617ba308e66512073d1fa72f19f399f00000000
