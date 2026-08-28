package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const Difficulty uint64 = 1

type Block struct {
	Index        uint64
	Timestamp    int64
	PreviousHash string
	Hash         string
	Nonce        uint64
	Difficulty   uint64
	Transactions []Transaction
}

type Blockchain struct {
	Blocks []Block
}

// ============================================================
// BLOCK CREATION
// ============================================================

func CreateBlock(
	index uint64,
	previousHash string,
	transactions []Transaction,
) Block {

	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Hash:         "",
		Nonce:        0,
		Difficulty:   Difficulty,
		Transactions: transactions,
	}

	block.Mine()

	return block
}

// ============================================================
// HASHING
// ============================================================

func (b Block) CalculateHash() string {

	data := strconv.FormatUint(
		b.Index,
		10,
	)

	data += strconv.FormatInt(
		b.Timestamp,
		10,
	)

	data += b.PreviousHash

	data += strconv.FormatUint(
		b.Nonce,
		10,
	)

	data += strconv.FormatUint(
		b.Difficulty,
		10,
	)

	for _, tx := range b.Transactions {
		data += tx.ID()
	}

	hash := sha256.Sum256(
		[]byte(data),
	)

	return hex.EncodeToString(
		hash[:],
	)
}

// ============================================================
// PROOF OF WORK
// ============================================================

func (b *Block) Mine() {

	target := strings.Repeat(
		"0",
		int(b.Difficulty),
	)

	for {

		b.Hash = b.CalculateHash()

		if strings.HasPrefix(
			b.Hash,
			target,
		) {
			return
		}

		b.Nonce++
	}
}

func (b Block) HasValidProofOfWork() bool {

	if b.Difficulty == 0 {
		return false
	}

	target := strings.Repeat(
		"0",
		int(b.Difficulty),
	)

	calculatedHash := b.CalculateHash()

	return calculatedHash == b.Hash &&
		strings.HasPrefix(
			b.Hash,
			target,
		)
}

// ============================================================
// BASIC BLOCK VALIDATION
// ============================================================

func ValidateBlock(
	block Block,
	previousBlock *Block,
) error {

	if block.Difficulty == 0 {
		return fmt.Errorf(
			"block %d has invalid difficulty",
			block.Index,
		)
	}

	if block.Hash == "" {
		return fmt.Errorf(
			"block %d has no hash",
			block.Index,
		)
	}

	if block.Index == 0 {

		if block.PreviousHash != "" {
			return fmt.Errorf(
				"genesis block has a previous hash",
			)
		}

	} else {

		if previousBlock == nil {
			return fmt.Errorf(
				"block %d is missing its previous block",
				block.Index,
			)
		}

		if block.Index != previousBlock.Index+1 {
			return fmt.Errorf(
				"invalid block index: got %d, expected %d",
				block.Index,
				previousBlock.Index+1,
			)
		}

		if block.PreviousHash != previousBlock.Hash {
			return fmt.Errorf(
				"block %d has an invalid previous hash",
				block.Index,
			)
		}

		if block.Timestamp < previousBlock.Timestamp {
			return fmt.Errorf(
				"block %d timestamp is earlier than previous block",
				block.Index,
			)
		}
	}

	if !block.HasValidProofOfWork() {
		return fmt.Errorf(
			"block %d has invalid proof-of-work",
			block.Index,
		)
	}

	return nil
}

// ============================================================
// TRANSACTION VALIDATION
// ============================================================

func ValidateBlockTransactions(
	block Block,
	utxoSet *UTXOSet,
) error {

	// Genesis block must contain no transactions.
	if block.Index == 0 {

		if len(block.Transactions) != 0 {
			return fmt.Errorf(
				"genesis block cannot contain transactions",
			)
		}

		return nil
	}

	if len(block.Transactions) == 0 {
		return fmt.Errorf(
			"block %d contains no transactions",
			block.Index,
		)
	}

	coinbaseCount := 0

	for i, tx := range block.Transactions {

		if tx.Coinbase {

			coinbaseCount++

			// Coinbase must be first.
			if i != 0 {
				return fmt.Errorf(
					"coinbase transaction must be first",
				)
			}

			err := ApplyCoinbaseTransaction(
				utxoSet,
				&tx,
				block.Index,
			)

			if err != nil {
				return fmt.Errorf(
					"invalid coinbase transaction: %w",
					err,
				)
			}

			continue
		}

		// Normal transactions must have inputs.
		if len(tx.Inputs) == 0 {
			return fmt.Errorf(
				"transaction %s has no inputs",
				tx.ID(),
			)
		}

		// Normal transactions must have outputs.
		if len(tx.Outputs) == 0 {
			return fmt.Errorf(
				"transaction %s has no outputs",
				tx.ID(),
			)
		}

		err := ApplyTransaction(
			utxoSet,
			&tx,
		)

		if err != nil {
			return fmt.Errorf(
				"invalid transaction %s: %w",
				tx.ID(),
				err,
			)
		}
	}

	if coinbaseCount != 1 {
		return fmt.Errorf(
			"block %d must contain exactly one coinbase transaction",
			block.Index,
		)
	}

	return nil
}

// ============================================================
// BLOCKCHAIN CREATION
// ============================================================

func CreateBlockchain() Blockchain {

	genesis := CreateBlock(
		0,
		"",
		[]Transaction{},
	)

	return Blockchain{
		Blocks: []Block{
			genesis,
		},
	}
}

// ============================================================
// ADD BLOCK
// ============================================================

func (bc *Blockchain) AddBlock(
	transactions []Transaction,
) {

	if len(bc.Blocks) == 0 {

		genesis := CreateBlock(
			0,
			"",
			[]Transaction{},
		)

		bc.Blocks = append(
			bc.Blocks,
			genesis,
		)
	}

	previousBlock := bc.Blocks[
		len(bc.Blocks)-1,
	]

	newBlock := CreateBlock(
		previousBlock.Index+1,
		previousBlock.Hash,
		transactions,
	)

	bc.Blocks = append(
		bc.Blocks,
		newBlock,
	)
}

// ============================================================
// FULL BLOCKCHAIN VALIDATION
// ============================================================

func IsBlockchainValid(
	bc Blockchain,
) bool {

	if len(bc.Blocks) == 0 {
		return false
	}

	for i, block := range bc.Blocks {

		if i == 0 {

			err := ValidateBlock(
				block,
				nil,
			)

			if err != nil {
				return false
			}

			continue
		}

		previousBlock := &bc.Blocks[i-1]

		err := ValidateBlock(
			block,
			previousBlock,
		)

		if err != nil {
			return false
		}
	}

	return true
}

// ============================================================
// PRINT BLOCK
// ============================================================

func PrintBlock(
	block Block,
) {

	fmt.Println(
		"Block:",
		block.Index,
	)

	fmt.Println(
		"Timestamp:",
		block.Timestamp,
	)

	fmt.Println(
		"Previous Hash:",
		block.PreviousHash,
	)

	fmt.Println(
		"Hash:",
		block.Hash,
	)

	fmt.Println(
		"Nonce:",
		block.Nonce,
	)

	fmt.Println(
		"Difficulty:",
		block.Difficulty,
	)

	fmt.Println(
		"Transactions:",
		len(block.Transactions),
	)

	fmt.Println(
		"Proof of Work:",
		block.HasValidProofOfWork(),
	)
}