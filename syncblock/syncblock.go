package syncblock

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gosuri/uilive"
)

const (
	redisDB        = 0
	blockInterval  = 10 * time.Minute
	maxRetries     = 5
	initialBackoff = 1 * time.Second
	capacity       = 2000
	// Blockstream Esplora API endpoint
	esploraBaseURL = "https://blockstream.info/api"
)

var ctx = context.Background()

type SyncBlock struct {
	blockstreamToken string
	network          string
	chain            string
	redisAddr        string
	redisPassword    string
}

type TransactionHistory struct {
	Height int    `json:"height"`
	TxHash string `json:"tx_hash"`
}

// Simplified Block structure
type Block struct {
	Height int      `json:"height"`
	Txids  []string `json:"txids"`
}
type Transaction struct {
	TxID     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
	Size     int    `json:"size"`
	Weight   int    `json:"weight"`
	Fee      int    `json:"fee"`
	Status   Status `json:"status"`
}

type Vin struct {
	TxID         string   `json:"txid"`
	Vout         int      `json:"vout"`
	Prevout      Prevout  `json:"prevout"`
	Scriptsig    string   `json:"scriptsig"`
	ScriptsigAsm string   `json:"scriptsig_asm"`
	Witness      []string `json:"witness"`
	IsCoinbase   bool     `json:"is_coinbase"`
	Sequence     uint32   `json:"sequence"`
}
type Prevout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

type Vout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

type Status struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}

var writer *uilive.Writer
var writer1 io.Writer
var writer2 io.Writer

func NewSync(blockstreamToken, network, chain, redisAddr, redisPassword string) *SyncBlock {

	return &SyncBlock{
		blockstreamToken: blockstreamToken,
		network:          network,
		chain:            chain,
		redisAddr:        redisAddr,
		redisPassword:    redisPassword,
	}
}

func (s *SyncBlock) Initiate() {
	client := redis.NewClient(&redis.Options{
		Addr:     s.redisAddr,
		Password: s.redisPassword,
		DB:       redisDB,
	})
	writer = uilive.New()
	writer1 = writer.Newline()
	writer2 = writer.Newline()

	writer.Start()

	defer writer.Stop()
	retries := 0
	backoff := initialBackoff

	for {
		err := processLatestBlocks(client)
		if err != nil {
			log.Printf("Error processing blocks: %v", err)
		}
		if isRateLimitError(err) {
			log.Printf("Rate limit hit, retrying in %v...", backoff)
			time.Sleep(backoff)
			backoff *= 2
			retries = -1
		}
		retries++
		time.Sleep(blockInterval)
	}
}

func processLatestBlocks(client *redis.Client) error {
	latestBlockHeight, err := getLatestBlockHeight()
	//fmt.Println("latestBlockHeight", latestBlockHeight)
	if err != nil {
		return fmt.Errorf("failed to get latest block height: %w", err)
	}
	//fmt.Fprintf(writer, "======= Syncing =======\n")
	for i := 0; i < 20; i++ {
		block, err := getBlockWithRetry(latestBlockHeight - i)

		if err != nil {
			log.Printf("[0] Failed to get block at height %d: %v", latestBlockHeight-i, err)
			continue
		}

		err = processBlock(block, client, i)
		if err != nil {
			log.Printf("Failed to process block at height %d: %v", latestBlockHeight-i, err)
		}
	}

	return nil
}

func getLatestBlockHeight() (int, error) {
	url := fmt.Sprintf("%s/blocks/tip/height", esploraBaseURL)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block height: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get latest block height, status code: %d", resp.StatusCode)
	}

	var height int
	err = json.NewDecoder(resp.Body).Decode(&height)
	if err != nil {
		return 0, fmt.Errorf("failed to decode latest block height response: %v", err)
	}

	return height, nil
}

func getBlockWithRetry(height int) (Block, error) {
	var block Block
	var err error
	backoff := initialBackoff

	for retries := 0; retries < maxRetries; retries++ {
		//fmt.Fprintf(writer1, "Block number: %d\n", i)
		block, err = getBlock(height)
		if err == nil {

			return block, nil
		}

		if isRateLimitError(err) {
			log.Printf("Rate limit hit when fetching block %d, retrying in %v...", height, backoff)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		break
	}

	return block, err
}

func getBlock(height int) (Block, error) {
	url := fmt.Sprintf("%s/block-height/%d", esploraBaseURL, height)

	//fmt.Println("Url", url)
	resp, err := http.Get(url)
	if err != nil {
		return Block{}, fmt.Errorf("[1] failed to get block at height  %d: %v", height, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("[2] failed to get block at height %d, status code: %d", height, resp.StatusCode)
	}

	var block Block
	var blockhash string

	body, _ := ioutil.ReadAll(resp.Body)
	blockhash = string(body)
	txid, _ := getBlockTxids(blockhash)

	block = Block{Height: height, Txids: txid}
	return block, nil
}
func getBlockTxids(hash string) ([]string, error) {
	url := fmt.Sprintf("%s/block/%s/txids", esploraBaseURL, hash)

	//fmt.Println("Url", url)
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, fmt.Errorf("[1] failed to get txid for blockhash  %s: %v", hash, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("[2] failed to get btxid for blockhash %s, status code: %d", hash, resp.StatusCode)
	}
	var blockhash string

	body, _ := ioutil.ReadAll(resp.Body)
	blockhash = string(body)
	var txids []string

	// Unmarshal the JSON array string into the array variable
	err = json.Unmarshal([]byte(blockhash), &txids)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return []string{}, fmt.Errorf("[2] failed to get btxid for blockhash %s, status code: %d", hash, resp.StatusCode)
	}
	return txids, nil
}

func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	// Check if error contains rate limit status code
	if strings.Contains(err.Error(), "429") {
		return true
	}
	return false
}

func processBlock(block Block, client *redis.Client, z int) error {
	var scrptIndx int = 0
	for i, txID := range block.Txids {
		// Simulated function for getting transaction details
		tx, err := getTransaction(txID)
		if err != nil {
			log.Printf("Failed to get transaction %s: %v", txID, err)
			continue
		}

		// Simulated processing of transaction outputs
		for j, out := range tx.Vout {

			//fmt.Println("Address", out.ScriptPubKeyAddress)
			if out.ScriptPubKeyAddress != "" {
				scripthash, err := getScriptHashFromAddress(out.ScriptPubKeyAddress)
				if err != nil {
					log.Printf("Failed to get scripthash for address %s: %v", out.ScriptPubKeyAddress, err)
					continue
				}
				scrptIndx += j
				history := TransactionHistory{
					Height: block.Height,
					TxHash: txID,
				}
				if scripthash != "" {

					err = UpdateRedis(client, scripthash, history)

					if err != nil {
						log.Printf("Failed to update Redis for scripthash %s: %v", scripthash, err)
					}
					fmt.Fprintf(writer, "======= Syncing =======\n")
					fmt.Fprintf(writer1, "Block number: %d\n", z)
					fmt.Fprintf(writer2, "Transaction number: (%d) Scripthash obtained and updated:%d\n", i, scrptIndx)
				}
			}
			//}
		}
	}

	return nil
}

/*
	func getScriptHashFromAddress(address string) (string, error) {
		addr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
		if err != nil {
			return "", err
		}
		scriptHash := btcutil.Hash160(addr.ScriptAddress())
		return hex.EncodeToString(scriptHash), nil
	}
*/
func getScriptHashFromAddress(address string) (string, error) {
	// Decode the Bitcoin address
	/*addr, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode address: %w", err)
	}

	// Convert the address to its script representation.
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", fmt.Errorf("failed to convert address to script: %w", err)
	}
	// Hash the scriptPubKey using SHA256
	hash := sha256.Sum256(script)
	// Convert the hash to a hexadecimal string and reverse it
	//reversedHash := reverseBytes(hash[:])*/

	return "", nil //hex.EncodeToString(hash[:]), nil
}

// reverseBytes reverses a byte slice.
func reverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
func UpdateRedis(client *redis.Client, scripthash string, history TransactionHistory) error {
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	if err := client.LPush(ctx, scripthash, historyJSON).Err(); err != nil {
		return err
	}

	// Trim the list to the specified capacity
	if err := client.LTrim(ctx, scripthash, 0, capacity-1).Err(); err != nil {
		return err
	}
	//log.Println("[*] Updated")
	return nil
}
func GetScripthashHistory(scripthash string) ([]TransactionHistory, error) {
	url := fmt.Sprintf("%s/scripthash/%s/txs", esploraBaseURL, scripthash)

	//fmt.Println("Url", url)
	resp, err := http.Get(url)
	if err != nil {
		return []TransactionHistory{}, fmt.Errorf("[1] failed to get ScriptHistory %s: %v", scripthash, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []TransactionHistory{}, fmt.Errorf("[2] failed to get ScriptHistory %s, status code: %d", scripthash, resp.StatusCode)
	}
	var trxhisty []TransactionHistory
	var blockhash string

	body, _ := ioutil.ReadAll(resp.Body)
	blockhash = string(body)
	//fmt.Println()
	var scriptH []Transaction
	if blockhash != "[]" {
		err = json.Unmarshal([]byte(blockhash), &scriptH)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return []TransactionHistory{}, fmt.Errorf("[2] failed to get ScriptHistory %s, status code: %d", scripthash, resp.StatusCode)
		}
		for _, tx := range scriptH {
			txn := TransactionHistory{Height: tx.Status.BlockHeight, TxHash: tx.TxID}
			trxhisty = append(trxhisty, txn)

		}
	} else {
		trxhisty = append(trxhisty, TransactionHistory{})
	}
	return trxhisty, nil
}

func getTransaction(txID string) (Transaction, error) {

	url := fmt.Sprintf("%s/tx/%s", esploraBaseURL, txID)

	//fmt.Println("Url", url)
	resp, err := http.Get(url)
	if err != nil {
		return Transaction{}, fmt.Errorf("[1] failed to get tx for txID  %s: %v", txID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Transaction{}, fmt.Errorf("[2] failed to get tx for txID %s, status code: %d", txID, resp.StatusCode)
	}
	var blockhash string

	body, _ := ioutil.ReadAll(resp.Body)
	blockhash = string(body)
	var tx Transaction

	// Unmarshal the JSON array string into the array variable
	err = json.Unmarshal([]byte(blockhash), &tx)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return Transaction{}, fmt.Errorf("[2] failed to get tx for txID %s, status code: %d", txID, resp.StatusCode)
	}
	//fmt.Println("Transaction", tx)
	return tx, nil
}
