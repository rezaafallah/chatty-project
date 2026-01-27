package crypto

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/tyler-smith/go-bip39"
)

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func HashString(input string) string {
	hash := sha512.Sum512([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateKeyPair is a mock. Use real Curve25519 in prod.
func GenerateKeyPair(mnemonic string) (string, string) {
	return "mock_pub_" + HashString(mnemonic)[:10], "mock_priv_key"
}