package handlers

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
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

	_, _ = fmt.Fprintf(channel, "Welcome to vis SSH server!\n\r")

	terminal := term.NewTerminal(channel, "> ")

	for {
		line, err := terminal.ReadLine()
		if err != nil {
			_, _ = fmt.Fprintf(channel, "error read line %e\n", err)
			break
		}

		if errH := handler(channel, line); errH != nil {
			log.Printf("error in handler %+v", errH)
		} else {
			_, _ = fmt.Fprintf(channel, "handle %s success\n\r", line)
		}
	}
	_, _ = fmt.Fprintf(channel, "closing session")

}

func handler(channel ssh.Channel, command string) error {
	switch command {
	case "tmux":
		return handleTmux(channel)
	case "exit":
		return handleExit(channel)
	default:
		_, _ = fmt.Fprintf(channel, "Unknown command: %s\n\r", command)
	}

	return nil
}

func handleTmux(channel ssh.Channel) error {
	_, _ = fmt.Fprintf(channel, "Upload binary to remoute server\n\r")
	if err := startTmux(channel); err != nil {
		return fmt.Errorf("failed to upload bin %w", err)
	}
	return nil
}

func startTmux(channel ssh.Channel) error {
	cmd := exec.Command("tmux")

	ptmx, errS := pty.Start(cmd)
	if errS != nil {
		return fmt.Errorf("failed to start tmux: %v", errS)
	}
	defer func(ptmx *os.File) {
		_ = ptmx.Close()
	}(ptmx)

	go func() {
		_, _ = io.Copy(channel, ptmx)
	}()

	go func() {
		_, _ = io.Copy(ptmx, channel)
	}()

	if errW := cmd.Wait(); errW != nil {
		return fmt.Errorf("tmux failed: %v", errW)
	}

	return nil
}

func handleExit(channel ssh.Channel) error {
	_, _ = fmt.Fprintf(channel, "Stoping server\n\r")
	os.Exit(0)
	return nil
}
