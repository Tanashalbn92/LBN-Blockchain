package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
)

type TxInput struct {
	TxID      string
	Index     uint32
	Signature []byte
}

type TxOutput struct {
	Owner  string
	Amount uint64
}

type Transaction struct {
	Inputs         []TxInput
	Outputs        []TxOutput
	Coinbase       bool
	CoinbaseHeight uint64
}

func (tx *Transaction) signingHash() []byte {
	copyTx := *tx

	copyTx.Inputs = make(
		[]TxInput,
		len(tx.Inputs),
	)

	for i, input := range tx.Inputs {
		copyTx.Inputs[i] = TxInput{
			TxID:  input.TxID,
			Index: input.Index,
		}
	}

	data, _ := json.Marshal(copyTx)

	hash := sha256.Sum256(data)

	return hash[:]
}

func (tx *Transaction) hash() []byte {
	data, _ := json.Marshal(tx)

	hash := sha256.Sum256(data)

	return hash[:]
}

func (tx *Transaction) ID() string {
	return hex.EncodeToString(
		tx.hash(),
	)
}

func (tx *Transaction) Sign(
	privateKey *ecdsa.PrivateKey,
) error {

	if tx.Coinbase {
		return errors.New(
			"coinbase transactions cannot be signed",
		)
	}

	if len(tx.Inputs) == 0 {
		return errors.New(
			"transaction has no inputs",
		)
	}

	hash := tx.signingHash()

	r, s, err := ecdsa.Sign(
		rand.Reader,
		privateKey,
		hash,
	)

	if err != nil {
		return err
	}

	signature := encodeSignature(
		r,
		s,
	)

	for i := range tx.Inputs {
		tx.Inputs[i].Signature = signature
	}

	return nil
}

func (tx *Transaction) VerifySignature(
	publicKey *ecdsa.PublicKey,
) bool {

	if tx.Coinbase {
		return false
	}

	if len(tx.Inputs) == 0 {
		return false
	}

	signature := tx.Inputs[0].Signature

	if len(signature) == 0 {
		return false
	}

	r, s, ok := decodeSignature(
		signature,
	)

	if !ok {
		return false
	}

	hash := tx.signingHash()

	return ecdsa.Verify(
		publicKey,
		hash,
		r,
		s,
	)
}

func encodeSignature(
	r *big.Int,
	s *big.Int,
) []byte {

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	result := make(
		[]byte,
		0,
		8+len(rBytes)+len(sBytes),
	)

	rLength := uint32(len(rBytes))
	sLength := uint32(len(sBytes))

	result = append(
		result,
		byte(rLength>>24),
		byte(rLength>>16),
		byte(rLength>>8),
		byte(rLength),
	)

	result = append(
		result,
		rBytes...,
	)

	result = append(
		result,
		byte(sLength>>24),
		byte(sLength>>16),
		byte(sLength>>8),
		byte(sLength),
	)

	result = append(
		result,
		sBytes...,
	)

	return result
}

func decodeSignature(
	signature []byte,
) (*big.Int, *big.Int, bool) {

	if len(signature) < 8 {
		return nil, nil, false
	}

	rLength := uint32(signature[0])<<24 |
		uint32(signature[1])<<16 |
		uint32(signature[2])<<8 |
		uint32(signature[3])

	rStart := 4
	rEnd := rStart + int(rLength)

	if rEnd+4 > len(signature) {
		return nil, nil, false
	}

	sLength := uint32(signature[rEnd])<<24 |
		uint32(signature[rEnd+1])<<16 |
		uint32(signature[rEnd+2])<<8 |
		uint32(signature[rEnd+3])

	sStart := rEnd + 4
	sEnd := sStart + int(sLength)

	if sEnd != len(signature) {
		return nil, nil, false
	}

	r := new(big.Int).SetBytes(
		signature[rStart:rEnd],
	)

	s := new(big.Int).SetBytes(
		signature[sStart:sEnd],
	)

	return r, s, true
}

func CreateCoinbaseTransaction(
	minerAddress string,
	reward uint64,
	blockHeight uint64,
) Transaction {

	return Transaction{
		Inputs: []TxInput{},

		Outputs: []TxOutput{
			{
				Owner:  minerAddress,
				Amount: reward,
			},
		},

		Coinbase:       true,
		CoinbaseHeight: blockHeight,
	}
}

func FormatLBN(
	amount uint64,
) string {

	whole := amount / 1_000_000_000

	return strconv.FormatUint(
		whole,
		10,
	) + " LBN"
}


// CreateTransaction creates a signed normal transaction.
//
// The sender's UTXOs are selected to cover the requested amount.
// Any remaining amount is returned to the sender as change.
func CreateTransaction(
	sender *Wallet,
	utxoSet *UTXOSet,
	recipient string,
	amount uint64,
) (Transaction, error) {

	if sender == nil {
		return Transaction{}, errors.New(
			"sender wallet is nil",
		)
	}

	if utxoSet == nil {
		return Transaction{}, errors.New(
			"UTXO set is nil",
		)
	}

	if recipient == "" {
		return Transaction{}, errors.New(
			"recipient address is empty",
		)
	}

	if amount == 0 {
		return Transaction{}, errors.New(
			"transaction amount must be greater than zero",
		)
	}

	selected, total, err := GetSpendableUTXOs(
		utxoSet,
		sender.Address,
		amount,
	)

	if err != nil {
		return Transaction{}, err
	}

	tx := Transaction{
		Inputs:  make([]TxInput, 0, len(selected)),
		Outputs: make([]TxOutput, 0, 2),
		Coinbase: false,
	}

	for _, utxo := range selected {
		tx.Inputs = append(
			tx.Inputs,
			TxInput{
				TxID:  utxo.TxID,
				Index: utxo.Index,
			},
		)
	}

	// Send the requested amount to the recipient.
	tx.Outputs = append(
		tx.Outputs,
		TxOutput{
			Owner:  recipient,
			Amount: amount,
		},
	)

	// Return any excess UTXO value to the sender.
	change := total - amount

	if change > 0 {
		tx.Outputs = append(
			tx.Outputs,
			TxOutput{
				Owner:  sender.Address,
				Amount: change,
			},
		)
	}

	// Sign the transaction with the sender's private key.
	if err := tx.Sign(sender.PrivateKey); err != nil {
		return Transaction{}, err
	}

	return tx, nil
}