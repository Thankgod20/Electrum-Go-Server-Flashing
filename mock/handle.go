package mock

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
type TransactionHistory struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

func (s *Server) handleRequest(message string, redisClient *redis.Client) string {
	var request JSONRPCRequest //map[string]interface{}
	var batchRequest []JSONRPCRequest
	err := json.Unmarshal([]byte(message), &request)
	if err == nil {
		return s.processSingleRequest(request, redisClient)
	}

	err = json.Unmarshal([]byte(message), &batchRequest)
	if err == nil {

		return s.processBatchRequest(batchRequest, redisClient)
	}
	log.Println("Invalid request:", err)
	return makeResponse(nil, nil, "Invalid request")
}
func (s *Server) processSingleRequest(request JSONRPCRequest, redisClient *redis.Client) string {
	switch request.Method {
	case "server.version":
		return makeResponse(request.ID, []interface{}{"ElectrumX 1.16.0", "1.4"}, nil)
	case "blockchain.headers.subscribe":
		//header, _ := s.bitcoinAPI.GetLatestBlockHeader()
		return handleGetMockHeader(request.ID, request.Params)
	case "blockchain.scripthash.get_history":
		return handleMockGetHistory(request.ID, request.Params) //s.handleGetHistory(request.ID, request.Params, redisClient)
	case "blockchain.scripthash.get_balance":
		return handleGetMockBalance(request.ID, request.Params)
	case "blockchain.scripthash.listunspent":
		return handleMockListUnspent(request.ID, request.Params)
	case "mempool.get_fee_histogram":
		return handleMockMempoolFee(request.ID)
	case "blockchain.estimatefee":
		return estimatedFee(request.ID, request.Params)
	case "blockchain.transaction.broadcast":
		return broadcastTransaction(request.ID, request.Params)
	case "blockchain.scripthash.get_mempool":
		return makeResponse(request.ID, nil, nil)
	case "blockchain.transaction.get":
		return makeResponse(request.ID, nil, nil)
	case "server.ping":
		return makeResponse(request.ID, nil, nil)
	default:
		return `{"error":"unknown method"}`
	}
}
func (s *Server) processBatchRequest(requests []JSONRPCRequest, redisClient *redis.Client) string {
	var responses []JSONRPCResponse
	for _, request := range requests {
		response := s.processSingleRequest(request, redisClient)
		var jsonResponse JSONRPCResponse
		json.Unmarshal([]byte(response), &jsonResponse)
		responses = append(responses, jsonResponse)
	}
	responseJSON, _ := json.Marshal(responses)
	return string(responseJSON)
}

func makeResponse(id interface{}, result interface{}, err interface{}) string {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
		Error:   err,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON)
}
