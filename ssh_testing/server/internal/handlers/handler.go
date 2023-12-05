package handlers

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func HandleSSHConnection(conn net.Conn, config *ssh.ServerConfig) {

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

		channel, _, errA := newChannel.Accept()
		if errA != nil {
			log.Printf("Failed to accept channel: %v", errA)
			return
		}

		go handleSSHSession(channel)
	}
}

func handleSSHSession(channel ssh.Channel) {
	defer func(channel ssh.Channel) {
		_ = channel.Close()
	}(channel)

	_, _ = fmt.Fprintf(channel, "Welcome to vis SSH server!")
	_, _ = fmt.Fprintf(channel, "\n\r")

	terminal := term.NewTerminal(channel, "> ")

	for {
		line, err := terminal.ReadLine()
		if err != nil {
			_, _ = fmt.Fprintf(channel, "error read line %e\n", err)
			break
		}

		response := handler(line)

		if response != "" {
			_, _ = fmt.Fprintf(channel, response)
			_, _ = fmt.Fprintf(channel, "\n\r")
		}
	}
}

// TODO: обрабатываю команды
func handler(command string) string {
	switch command {
	case "exit":
		return "Goodbye!"
	case "start vis":
		return "Starting vis..."
	case "push bin":
		return "Pushing binary..."
	default:
		return "Unknown command: " + command
	}
}
