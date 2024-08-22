package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
	"regexp"
	"strings"

	"github.com/Thankgod20/miniBTCD/blockchain"
)

type BitcoinRPCClient struct {
	Client *rpc.Client
}

func NewBitcoinRPCClient(rpcurl string) *BitcoinRPCClient {
	client, err := rpc.Dial("tcp", rpcurl) //"localhost:18885"
	if err != nil {
		log.Fatalf("Failed to connect to RPC server: %v", err)
	}
	log.Println("Conneceted")
	return &BitcoinRPCClient{Client: client}
}
func (bc *BitcoinRPCClient) GetBlockHeader() (string, string) {
	args := blockchain.GetLatestBlockArgs{}
	var reply blockchain.GetLatestBlockReply

	err := bc.Client.Call("Blockchain.GetLatestBlock", &args, &reply)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}
	// Prepare the data by removing newline characters and properly formatting
	re := regexp.MustCompile(`(?m)^\s*`)
	reply.JSONString = re.ReplaceAllString(reply.JSONString, "")

	// Split the data into lines
	lines := strings.Split(reply.JSONString, "\n")

	// Prepare a map to hold key-value pairs
	dataMap := make(map[string]string)

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split each line into key and value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			dataMap[key] = value
		}
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", "" //-1
	}

	// Unmarshal JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return "", "" // -1
	}

	// Extract the hex value
	hexValue, exists := result["Hex"]

	if exists {

		return hexValue.(string), result["Height"].(string) //height //strconv.ParseInt(result["Height"].(string), 10, -1)
	} else {

		return "", "" //-1
	}

}
func (bc *BitcoinRPCClient) GetBalance(address string) string {
	args := blockchain.GetBalanceArgs{Address: address}
	var reply blockchain.GetLatestBlockReply
	err := bc.Client.Call("Blockchain.GetBalanceByHash", &args, &reply)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}

	var balance int
	err = json.Unmarshal(reply.JSONBlock, &balance)
	if err != nil {
		log.Fatalf("Failed to unmarshal block JSON: %v", err)
	}
	balData := fmt.Sprintf("{\"confirmed\":%d,\"unconfirmed\":\"%d\"}", balance, 0)
	return balData
}
func (bc *BitcoinRPCClient) GetTransactionHistory(scripthash string) ([]map[string]interface{}, []map[string]interface{}) {
	args := blockchain.GetAddressHistoryArgs{Address: scripthash}
	var reply blockchain.GetAddressHistoryReply

	err := bc.Client.Call("Blockchain.GetTransactionHistoryScriptHash", &args, &reply)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}
	var trxall []map[string]interface{}
	var mempooltrxall []map[string]interface{}
	for _, txId := range reply.TransactionHex {
		argsx := blockchain.GetTransactionReply{TransactionID: txId}
		var replyx blockchain.GetLatestBlockReply
		err := bc.Client.Call("Blockchain.GetFulTX", &argsx, &replyx)
		if err != nil {
			log.Fatalf("Failed to get latest block: %v", err)
		}
		result := parseJson(replyx.JSONString)
		//log.Println("Resusjttlt", result, replyx.JSONString)
		height := result["Height"].(string)

		balData := fmt.Sprintf("{\"height\":%s,\"tx_hash\":\"%s\"}", height, txId)
		var trx map[string]interface{}
		err = json.Unmarshal([]byte(balData), &trx)
		if err != nil {
			fmt.Printf("Failed to Parse Balance data: %v", err)
		}
		trxall = append(trxall, trx)
		//log.Println("Block Heig", result["Height"])
	}
	for _, txid := range reply.TransactionHexMempool {
		balData := fmt.Sprintf("{\"fee\":%d,\"height\":%d,\"tx_hash\":\"%s\"}", 2000, 0, txid)
		var trx map[string]interface{}
		err = json.Unmarshal([]byte(balData), &trx)
		if err != nil {
			fmt.Printf("Failed to Parse Balance data: %v", err)
		}
		trxall = append(trxall, trx)
		mempooltrxall = append(mempooltrxall, trx)
	}
	//log.Println("Transactions:", trxall, mempooltrxall)
	//
	return trxall, mempooltrxall
}
func (bc *BitcoinRPCClient) GetTransactionUTXOs(scripthash string) []map[string]interface{} {
	args := blockchain.GetAddressHistoryArgs{Address: scripthash}
	var reply blockchain.GetAddressHistoryReply

	err := bc.Client.Call("Blockchain.GetUTXOSScripttHash", &args, &reply)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}
	var allutxos []map[string]interface{}

	for _, utxos := range reply.TransactionHex {
		var result map[string]interface{}
		err := json.Unmarshal([]byte(utxos), &result)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		allutxos = append(allutxos, result)
		//log.Println("Block Heig", result["Height"])
	}

	//log.Println("Transactions:", trxall, mempooltrxall)
	//
	return allutxos
}
func (bc *BitcoinRPCClient) GetTransaction(txId string) (string, map[string]interface{}) {
	argsx := blockchain.GetTransactionReply{TransactionID: txId}
	var replyx blockchain.GetLatestBlockReply
	err := bc.Client.Call("Blockchain.GetFulTXElect", &argsx, &replyx)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}
	var result map[string]interface{}
	log.Println("TRXID:", txId)
	log.Println("JSONS:", replyx.JSONString)
	if replyx.JSONString == "" {
		replyx.JSONString = "{}"
	}
	// Parse the JSON string into the map
	err = json.Unmarshal([]byte(replyx.JSONString), &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	hex := result["transactionHex"].(string)
	delete(result, "transactionHex")
	return hex, result
}
func (bc *BitcoinRPCClient) BroadCastTX(trx string) string {
	args := blockchain.GetTransactionsArgs{TransactionHex: trx}
	var reply blockchain.GetLatestBlockReply
	err := bc.Client.Call("Blockchain.AddToMempool", &args, &reply)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}
	var trnxID string
	if reply.JSONString == "" {
		trnxID = hex.EncodeToString(reply.JSONBlock)
		fmt.Println("Transaction ID", trnxID)
	} else {
		fmt.Println("Error", reply.JSONString)
	}
	return trnxID
}
func parseJson(jsons string) map[string]interface{} {
	re := regexp.MustCompile(`(?m)^\s*`)
	jsons = re.ReplaceAllString(jsons, "")

	// Split the data into lines
	lines := strings.Split(jsons, "\n")

	// Prepare a map to hold key-value pairs
	dataMap := make(map[string]string)

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split each line into key and value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			dataMap[key] = value
		}
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}

	// Unmarshal JSON into a map
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil // -1
	}
	return result
}
