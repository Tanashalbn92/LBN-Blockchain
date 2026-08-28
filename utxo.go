package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type UTXO struct {
	TxID   string
	Index  uint32
	Owner  string
	Amount uint64
}

type UTXOSet struct {
	Outputs map[string]UTXO
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		Outputs: make(map[string]UTXO),
	}
}

func UTXOKey(txID string, index uint32) string {
	return txID + ":" + strconv.FormatUint(uint64(index), 10)
}

func CreateUTXO(txID string, index uint32, owner string, amount uint64) UTXO {
	return UTXO{
		TxID:   txID,
		Index:  index,
		Owner:  owner,
		Amount: amount,
	}
}

func CalculateUTXOSetID(set *UTXOSet) string {
	data := ""

	for key, utxo := range set.Outputs {
		data += key
		data += utxo.Owner
		data += strconv.FormatUint(utxo.Amount, 10)
	}

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func (set *UTXOSet) Add(utxo UTXO) {
	key := UTXOKey(utxo.TxID, utxo.Index)
	set.Outputs[key] = utxo
}

func (set *UTXOSet) Remove(txID string, index uint32) {
	key := UTXOKey(txID, index)
	delete(set.Outputs, key)
}

func (set *UTXOSet) Balance(address string) uint64 {
	var balance uint64

	for _, utxo := range set.Outputs {
		if utxo.Owner == address {
			balance += utxo.Amount
		}
	}

	return balance
}