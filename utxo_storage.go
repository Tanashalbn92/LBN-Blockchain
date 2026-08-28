package main

import (
	"encoding/json"
	"errors"
	"os"
)

const utxoFile = "utxo.json"

func SaveUTXOSet(set *UTXOSet) error {
	data, err := json.MarshalIndent(set, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(
		utxoFile,
		data,
		0644,
	)
}

func LoadUTXOSet() (*UTXOSet, error) {
	data, err := os.ReadFile(utxoFile)

	if err != nil {
		return nil, err
	}

	var set UTXOSet

	err = json.Unmarshal(data, &set)

	if err != nil {
		return nil, err
	}

	if set.Outputs == nil {
		return nil, errors.New(
			"UTXO set contains no outputs",
		)
	}

	return &set, nil
}

func UTXOSetExists() bool {
	_, err := os.Stat(utxoFile)

	return err == nil
}