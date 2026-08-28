package main

import "errors"

const MaxSupply uint64 = 1_000_000_000

const InitialBlockReward uint64 = 50_000_000_000

const AtomicUnitsPerLBN uint64 = 1_000_000_000

const HalvingInterval uint64 = 100_000

// MaxSupplyAtomicUnits returns the maximum supply
// in the smallest LBN denomination.
func MaxSupplyAtomicUnits() uint64 {
	return MaxSupply * AtomicUnitsPerLBN
}

// BlockReward returns the allowed mining reward
// for a particular block height.
func BlockReward(height uint64) uint64 {

	if height == 0 {
		return 0
	}

	reward := InitialBlockReward

	halvings := (height - 1) / HalvingInterval

	for i := uint64(0); i < halvings; i++ {

		reward /= 2

		if reward == 0 {
			return 0
		}
	}

	return reward
}

// GetSpendableUTXOs selects UTXOs belonging to an
// address until the requested amount is covered.
func GetSpendableUTXOs(
	set *UTXOSet,
	address string,
	amount uint64,
) ([]UTXO, uint64, error) {

	var selected []UTXO
	var total uint64

	for _, utxo := range set.Outputs {

		if utxo.Owner != address {
			continue
		}

		if ^uint64(0)-total < utxo.Amount {
			return nil, 0, errors.New(
				"balance overflow",
			)
		}

		selected = append(
			selected,
			utxo,
		)

		total += utxo.Amount

		if total >= amount {
			return selected, total, nil
		}
	}

	return nil, 0, errors.New(
		"insufficient balance",
	)
}

// ApplyTransaction applies a normal transaction
// to the UTXO set.
func ApplyTransaction(
	set *UTXOSet,
	tx *Transaction,
) error {

	if tx.Coinbase {
		return errors.New(
			"coinbase requires block-height validation",
		)
	}

	if len(tx.Inputs) == 0 {
		return errors.New(
			"transaction has no inputs",
		)
	}

	if len(tx.Outputs) == 0 {
		return errors.New(
			"transaction has no outputs",
		)
	}

	var inputTotal uint64
	var outputTotal uint64

	seenInputs := make(
		map[string]bool,
	)

	for _, input := range tx.Inputs {

		key := UTXOKey(
			input.TxID,
			input.Index,
		)

		if seenInputs[key] {
			return errors.New(
				"double-spend detected",
			)
		}

		seenInputs[key] = true

		utxo, exists := set.Outputs[key]

		if !exists {
			return errors.New(
				"input UTXO does not exist",
			)
		}

		if len(input.Signature) == 0 {
			return errors.New(
				"input has no signature",
			)
		}

		if ^uint64(0)-inputTotal < utxo.Amount {
			return errors.New(
				"input amount overflow",
			)
		}

		inputTotal += utxo.Amount
	}

	for _, output := range tx.Outputs {

		if output.Owner == "" {
			return errors.New(
				"output has no owner",
			)
		}

		if output.Amount == 0 {
			return errors.New(
				"output amount must be greater than zero",
			)
		}

		if ^uint64(0)-outputTotal < output.Amount {
			return errors.New(
				"output amount overflow",
			)
		}

		outputTotal += output.Amount
	}

	if outputTotal > inputTotal {
		return errors.New(
			"transaction spends more than its inputs",
		)
	}

	for _, input := range tx.Inputs {

		set.Remove(
			input.TxID,
			input.Index,
		)
	}

	for index, output := range tx.Outputs {

		utxo := CreateUTXO(
			tx.ID(),
			uint32(index),
			output.Owner,
			output.Amount,
		)

		set.Add(utxo)
	}

	return nil
}

// ApplyCoinbaseTransaction applies a mining reward
// after verifying the reward against block height.
func ApplyCoinbaseTransaction(
	set *UTXOSet,
	tx *Transaction,
	height uint64,
) error {

	if !tx.Coinbase {
		return errors.New(
			"not a coinbase transaction",
		)
	}

	if len(tx.Inputs) != 0 {
		return errors.New(
			"coinbase cannot have inputs",
		)
	}

	if len(tx.Outputs) != 1 {
		return errors.New(
			"coinbase must have exactly one output",
		)
	}

	output := tx.Outputs[0]

	if output.Owner == "" {
		return errors.New(
			"coinbase has no owner",
		)
	}

	expectedReward := BlockReward(height)

	if expectedReward == 0 {
		return errors.New(
			"no reward available at this height",
		)
	}

	if output.Amount != expectedReward {
		return errors.New(
			"invalid block reward",
		)
	}

	currentSupply := CalculateSupply(set)

	maxSupply := MaxSupplyAtomicUnits()

	if currentSupply > maxSupply {
		return errors.New(
			"current supply exceeds maximum supply",
		)
	}

	if output.Amount > maxSupply-currentSupply {
		return errors.New(
			"block reward would exceed maximum supply",
		)
	}

	utxo := CreateUTXO(
		tx.ID(),
		0,
		output.Owner,
		output.Amount,
	)

	set.Add(utxo)

	return nil
}

// CalculateSupply calculates the total supply
// represented by the current UTXO set.
func CalculateSupply(
	set *UTXOSet,
) uint64 {

	var total uint64

	for _, utxo := range set.Outputs {

		if ^uint64(0)-total < utxo.Amount {
			return ^uint64(0)
		}

		total += utxo.Amount
	}

	return total
}

// ============================================================
// REBUILD UTXO SET
// ============================================================

// RebuildUTXOSet reconstructs the UTXO set by replaying
// every transaction in the blockchain from beginning to end.
//
// This allows a node to derive balances from the blockchain
// instead of trusting a separately stored balance database.
func RebuildUTXOSet(
	bc Blockchain,
) (*UTXOSet, error) {

	set := NewUTXOSet()

	if len(bc.Blocks) == 0 {
		return set, errors.New(
			"cannot rebuild UTXO set from empty blockchain",
		)
	}

	for blockIndex, block := range bc.Blocks {

		// Genesis block has no transactions.
		if blockIndex == 0 {

			if len(block.Transactions) != 0 {
				return nil, errors.New(
					"genesis block contains transactions",
				)
			}

			continue
		}

		if len(block.Transactions) == 0 {
			return nil, errors.New(
				"block contains no transactions",
			)
		}

		coinbaseCount := 0

		for txIndex := range block.Transactions {

			tx := block.Transactions[txIndex]

			if tx.Coinbase {

				coinbaseCount++

				if txIndex != 0 {
					return nil, errors.New(
						"coinbase transaction must be first",
					)
				}

				err := ApplyCoinbaseTransaction(
					set,
					&tx,
					block.Index,
				)

				if err != nil {
					return nil, err
				}

				continue
			}

			err := ApplyTransaction(
				set,
				&tx,
			)

			if err != nil {
				return nil, err
			}
		}

		if coinbaseCount != 1 {
			return nil, errors.New(
				"block must contain exactly one coinbase transaction",
			)
		}
	}

	return set, nil
}