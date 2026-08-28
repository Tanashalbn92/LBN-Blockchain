package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

const minerWalletFile = "miner_wallet.key"
const recipientWalletFile = "recipient_wallet.key"

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

// CreateWallet creates a new LBN wallet.
func CreateWallet() (*Wallet, error) {

	privateKey, err := ecdsa.GenerateKey(
		elliptic.P256(),
		rand.Reader,
	)

	if err != nil {
		return nil, err
	}

	return walletFromPrivateKey(privateKey)
}

// walletFromPrivateKey builds the public key and address
// from an existing private key.
func walletFromPrivateKey(
	privateKey *ecdsa.PrivateKey,
) (*Wallet, error) {

	if privateKey == nil {
		return nil, errors.New(
			"private key is nil",
		)
	}

	publicKey := append(
		privateKey.PublicKey.X.Bytes(),
		privateKey.PublicKey.Y.Bytes()...,
	)

	hash := sha256.Sum256(publicKey)

	address := hex.EncodeToString(hash[:])

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// SaveWallet permanently saves the wallet's private key.
func SaveWallet(
	wallet *Wallet,
	filename string,
) error {

	if wallet == nil || wallet.PrivateKey == nil {
		return errors.New(
			"wallet or private key is nil",
		)
	}

	data, err := x509.MarshalECPrivateKey(
		wallet.PrivateKey,
	)

	if err != nil {
		return err
	}

	encoded := hex.EncodeToString(data)

	return os.WriteFile(
		filename,
		[]byte(encoded),
		0600,
	)
}

// LoadWallet loads a previously saved wallet.
func LoadWallet(
	filename string,
) (*Wallet, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	der, err := hex.DecodeString(
		string(data),
	)

	if err != nil {
		return nil, errors.New(
			"invalid wallet file",
		)
	}

	privateKey, err := x509.ParseECPrivateKey(
		der,
	)

	if err != nil {
		return nil, err
	}

	return walletFromPrivateKey(
		privateKey,
	)
}

// LoadOrCreateWallet loads an existing wallet.
// If one does not exist, it creates and saves a new one.
func LoadOrCreateWallet(
	filename string,
) (*Wallet, error) {

	if _, err := os.Stat(filename); err == nil {

		return LoadWallet(filename)
	}

	wallet, err := CreateWallet()

	if err != nil {
		return nil, err
	}

	if err := SaveWallet(
		wallet,
		filename,
	); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (w *Wallet) Display() {

	fmt.Println(
		"===== LBN WALLET =====",
	)

	fmt.Println(
		"Address:",
		w.Address,
	)

	fmt.Println(
		"Public Key:",
		hex.EncodeToString(w.PublicKey),
	)
}