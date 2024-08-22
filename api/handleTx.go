package api

import (
	"context"
	"electrum-server-go/syncblock"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func (s *Server) handleGetHistory(id interface{}, params []interface{}, redisClient *redis.Client) string {
	if len(params) != 1 {
		return makeResponse(id, nil, "invalid params")
	}
	scripthash, ok := params[0].(string)
	if !ok {
		return makeResponse(id, nil, "invalid scripthash")
	}
	scripthashr, _ := reverseHexString(scripthash)
	history, err := queryRedis(redisClient, scripthashr)
	if err != nil {
		log.Println("Unable to Query Redis:", err)

	}

	// Not in database get
	if len(history) == 0 {

		trxHistyArray, err := syncblock.GetScripthashHistory(scripthashr)
		if err != nil {
			log.Println("Unable to get Transaction History:", err)
		}
		//fmt.Println("Index Legnt", trxHistyArray)
		if len(trxHistyArray) > 0 {
			syncblock.UpdateRedis(redisClient, scripthash, trxHistyArray[0])
			history, _ = queryRedis(redisClient, scripthash)
		}
	}
	//check height
	var hist_response string
	if len(history) > 0 {
		heights := history[0].(map[string]interface{})["height"].(float64)

		if heights == 0 {
			hist_response = makeResponse(id, []TransactionHistory{}, nil)
		} else {
			hist_response = makeResponse(id, history, nil)
		}
	} else {
		hist_response = makeResponse(id, []TransactionHistory{}, nil)
	}
	//fmt.Println("Trx History:", heights)

	return hist_response

}

// reverseHexString reverses the bytes of a hex string.
func reverseHexString(hexStr string) (string, error) {
	// Decode the hex string to bytes
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Reverse the bytes
	reversedBytes := reverseBytes(bytes)

	// Encode the reversed bytes back to a hex string
	reversedHexStr := hex.EncodeToString(reversedBytes)

	return reversedHexStr, nil
}
func reverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
func queryRedis(client *redis.Client, scripthash string) ([]interface{}, error) {
	historyJSONs, err := client.LRange(context.Background(), scripthash, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var history []interface{}
	for _, historyJSON := range historyJSONs {
		var txHistory map[string]interface{}
		err := json.Unmarshal([]byte(historyJSON), &txHistory)
		if err != nil {
			return nil, err
		}
		history = append(history, txHistory)
	}
	return history, nil
}
