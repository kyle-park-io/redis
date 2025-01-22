package usertest

import (
	"bufio"
	"net"
	"os"
	"testing"
	"time"

	"kyle-redis/client"
	"kyle-redis/handler"
	"kyle-redis/internal/config"
	"kyle-redis/logger"

	"github.com/spf13/viper"
)

// go test -v -run TestHandleClientConnection2
func TestHandleClientConnection2(t *testing.T) {

	env := "dev"
	// config
	config.SetEnv(env)
	client.Init()

	port := viper.GetString("port.server")
	if port == "" {
		logger.Log.Fatalln("Check env: port.server")
	}

	// 1. Start the server by listening on a specific port.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Log.Fatalf("Error starting server: %+v", err)
	}
	defer listener.Close()
	logger.Log.Infof("Server is listening on port %s...", port)

	// tcp client
	go simulateClientConnection2()

	for {
		// 2. Wait for client connections.
		conn, err := listener.Accept()
		if err != nil {
			logger.Log.Errorf("Error accepting connection: %+v", err)
			continue
		}
		logger.Log.Infoln("Client connected: ", conn.RemoteAddr())

		// 3. Handle each client connection in a separate goroutine.
		go handler.RedisHandleConnection(conn)
	}
}

func simulateClientConnection2() {
	// 1. Connect to the server.
	conn, err := net.Dial("tcp", ":"+viper.GetString("port.server"))
	if err != nil {
		logger.Log.Fatalf("Error connecting to server: %+v", err)
	}
	defer conn.Close()
	logger.Log.Infoln("Connected to the server. Type commands (e.g. PING, HELLO, TIME, EXIT)")

	for {
		conn.Write([]byte("REAL-TIME" + "\n"))
		// 3. Receive the response from the server.
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Log.Errorf("Error reading response: %+v", err)
			return
		}
		logger.Log.Infof("Response: %s", response)

		time.Sleep(2 * time.Second)
	}

	os.Exit(1)
}
