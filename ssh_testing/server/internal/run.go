package internal

import (
	"log"
	"net"

	"github.com/VyacheslavIsWorkingNow/siv/ssh_testing/server/internal/handlers"
)

func RunServer() {
	pk, errPK := parsePrivateKey(privateKeyPath)
	if errPK != nil {
		log.Fatal("Failed to parse private key:", errPK)
	}

	config := NewConfig()

	config.AddHostKey(pk)

	listener, errL := net.Listen("tcp", "0.0.0.0:2222")
	if errL != nil {
		log.Fatal("Failed to listen:", errL)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	log.Println("SSH server listening on :2222...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handlers.HandleSSHConnection(conn, config)
	}
}
