package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
)

func handleSSHConnection(conn net.Conn, config *ssh.ServerConfig) {

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer func(sshConn *ssh.ServerConn) {
		_ = sshConn.Close()
	}(sshConn)

	log.Printf("New SSH connection from %s", sshConn.RemoteAddr())

	go ssh.DiscardRequests(reqs)
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			err = newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			if err != nil {
				return
			}
			continue
		}

		channel, _, err := newChannel.Accept()
		if err != nil {
			log.Printf("Failed to accept channel: %v", err)
			return
		}

		go handleSSHSession(channel)
	}
}

func handleSSHSession(channel ssh.Channel) {
	defer func(channel ssh.Channel) {
		_ = channel.Close()
	}(channel)
	// TODO: добавить логику работы на сервере -> запуск бинарника
	_, err := fmt.Fprint(channel, "Welcome to my simple SSH server!\n")
	if err != nil {
		return
	}
}
