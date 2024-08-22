package mock

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
)

type Server struct {
	address       string
	redisAddr     string
	redisPassWord string
	//bitcoinAPI    *BitcoinAPIClient
	certFile string
	keyFile  string
}

func NewServer(address, redisAddr, redisPassWord, certFile, keyFile string) (*Server, error) {

	//bitcoinClientAPI := NewBitcoinAPIClient(apiurl, apiKey)
	return &Server{address: address, redisAddr: redisAddr, redisPassWord: redisPassWord, certFile: certFile, keyFile: keyFile}, nil
}

func (s *Server) Start() error {
	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", s.address, config) //net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     s.redisAddr,
		Password: s.redisPassWord,
		DB:       0,
	})
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error Accepting connection: ", err)
			continue
		}
		go s.handleConnection(conn, redisClient)
	}
}
func (s *Server) handleConnection(conn net.Conn, redisClient *redis.Client) {
	defer conn.Close()
	log.Println("New client connected:", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println("Client disconnected", conn.RemoteAddr().String())
			return
		}
		if err != nil {
			log.Println("Error reading from connection: ", err)
			return
		}
		//getblock, _ := s.bitcoinAPI.GetLatestBlockHeader()
		log.Println("Received message:", message)
		response := s.handleRequest(message, redisClient)
		log.Println("Sending response:", response)
		conn.Write([]byte(response + "\n"))
	}
}
