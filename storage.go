package main

import (
	"encoding/json"
	"errors"
	"os"
)

const blockchainFile = "blockchain.json"

func SaveBlockchain(bc Blockchain) error {
	data, err := json.MarshalIndent(
		bc,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		blockchainFile,
		data,
		0644,
	)
}

func LoadBlockchain() (Blockchain, error) {

	data, err := os.ReadFile(
		blockchainFile,
	)

	if err != nil {
		return Blockchain{}, err
	}

	var blockchain Blockchain

	err = json.Unmarshal(
		data,
		&blockchain,
	)

	if err != nil {
		return Blockchain{}, err
	}

	if len(blockchain.Blocks) == 0 {
		return Blockchain{}, errors.New(
			"blockchain contains no blocks",
		)
	}

	return blockchain, nil
}

func BlockchainExists() bool {
	_, err := os.Stat(blockchainFile)

	return err == nil
}