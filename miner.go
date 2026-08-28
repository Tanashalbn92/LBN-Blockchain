package main

import (
	"errors"
	"time"
)

// MineBlock creates a new block containing a coinbase
// transaction that pays the miner the correct block reward.
func MineBlock(
	bc *Blockchain,
	minerAddress string,
	transactions []Transaction,
) (Block, error) {

	if bc == nil {
		return Block{}, errors.New("blockchain is nil")
	}

	if minerAddress == "" {
		return Block{}, errors.New("miner address is empty")
	}

	if len(bc.Blocks) == 0 {
		return Block{}, errors.New("blockchain has no genesis block")
	}

	// The next block's height.
	blockHeight := uint64(len(bc.Blocks))

	// Get the scheduled mining reward.
	reward := BlockReward(blockHeight)

	if reward == 0 {
		return Block{}, errors.New("block reward is zero")
	}

	// Create a unique coinbase transaction for this
	// specific block height.
	coinbase := CreateCoinbaseTransaction(
		minerAddress,
		reward,
		blockHeight,
	)

	// Coinbase must be the first transaction.
	allTransactions := make(
		[]Transaction,
		0,
		len(transactions)+1,
	)

	allTransactions = append(
		allTransactions,
		coinbase,
	)

	allTransactions = append(
		allTransactions,
		transactions...,
	)

	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	// Build the new block.
	block := Block{
		Index:        blockHeight,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousBlock.Hash,
		Difficulty:   previousBlock.Difficulty,
		Transactions: allTransactions,
	}

	// Perform Proof of Work.
	block.Mine()

	return block, nil
}