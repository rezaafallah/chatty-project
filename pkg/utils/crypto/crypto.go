package crypto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"

	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic 
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// HashString
func HashString(input string) string {
	hash := sha512.Sum512([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateKeyPair
func GenerateKeyPair(mnemonic string) (publicKey string, privateKey string) {
	seed := bip39.NewSeed(mnemonic, "")
	reader := bytes.NewReader(seed)
	
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return "", ""
	}

	return hex.EncodeToString(pub), hex.EncodeToString(priv)
}