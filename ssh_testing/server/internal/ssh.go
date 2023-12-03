package internal

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	privateBytes, errR := os.ReadFile(keyPath)
	if errR != nil {
		return nil, fmt.Errorf("failed read ssh key %w", errR)
	}

	pk, errPPK := ssh.ParsePrivateKey(privateBytes)
	if errPPK != nil {
		return nil, fmt.Errorf("failed parse ssh key %w", errPPK)
	}

	return pk, nil
}
