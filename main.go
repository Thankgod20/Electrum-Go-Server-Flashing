package main

import (
	"electrum-server-go/api"
	"electrum-server-go/mock"
	"electrum-server-go/rpc"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) > 2 {
		var BTCAPI, BTCAPIKEY string
		var RPCURL string
		if strings.Contains(args[1], "--api=") {
			_btcapi := strings.Split(args[1], "=")
			BTCAPI = _btcapi[1]

			if strings.Contains(args[2], "--key=") {
				_btcapikey := strings.Split(args[2], "=")
				BTCAPIKEY = _btcapikey[1]
			}
			log.Println("API DETAILS", BTCAPI, BTCAPIKEY)
			go func() {
				//syncer := syncblock.NewSync("77608741d9074726bfca10d512936a89", "btc", "main", "127.0.0.1:6379", "")
				//syncer.Initiate()
			}()

			server, err := api.NewServer(":50001", "127.0.0.1:6379", "", BTCAPI, BTCAPIKEY, "cert.pem", "key.pem")
			if err != nil {
				log.Fatalf("Error creating server: %v", err)
			}
			log.Println("Starting Electrum server for btc node api on :50001")
			if err := server.Start(); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		} else if strings.Contains(args[1], "--rpc") {
			//_btcapi := strings.Split(args[1], "=")
			//BTCAPI = _btcapi[1]

			if strings.Contains(args[2], "--url=") {
				_btcapikey := strings.Split(args[2], "=")
				RPCURL = _btcapikey[1]
			}
			log.Println("RPC DETAILS", RPCURL)

			server, err := rpc.NewServer(":50001", "127.0.0.1:6379", "", RPCURL, "cert.pem", "key.pem")
			if err != nil {
				log.Fatalf("Error creating server: %v", err)
			}
			log.Println("Starting Electrum server for btc node api on :50001")
			if err := server.Start(); err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
		}
	} else if len(args) == 2 {
		log.Println("\n* Incomplete Parameters\n - ./electrumServer --api=https://api.com --key=SCJ...")
	} else {
		log.Println("Starting Mock Server")
		server, err := mock.NewServer(":50001", "http://localhost:8332", "", "cert.pem", "key.pem")
		if err != nil {
			log.Fatalf("Error creating server: %v", err)
		}
		log.Println("Starting Electrum server on :50001")
		if err := server.Start(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}
