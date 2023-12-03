package internal

import "golang.org/x/crypto/ssh"

const privateKeyPath = "/Users/slavaruswarrior/Documents/GitHub/siv/ssh_testing/private_key.pem"

func NewConfig() *ssh.ServerConfig {
	return &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}
}
