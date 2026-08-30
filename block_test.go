package main

import "testing"

func TestValidBlockchain(t *testing.T) {

	bc := CreateBlockchain()

	bc.AddBlock([]Transaction{})

	if !IsBlockchainValid(bc) {
		t.Fatal("valid blockchain was rejected")
	}
}

func TestTamperedBlockData(t *testing.T) {

	bc := CreateBlockchain()

	bc.AddBlock([]Transaction{})

	bc.Blocks[1].Timestamp++

	if IsBlockchainValid(bc) {
		t.Fatal("tampered timestamp was accepted")
	}
}

func TestTamperedBlockHash(t *testing.T) {

	bc := CreateBlockchain()

	bc.AddBlock([]Transaction{})

	bc.Blocks[1].Hash = "fake-hash"

	if IsBlockchainValid(bc) {
		t.Fatal("tampered hash was accepted")
	}
}

func TestBrokenPreviousHash(t *testing.T) {

	bc := CreateBlockchain()

	bc.AddBlock([]Transaction{})

	bc.Blocks[1].PreviousHash = "fake-previous-hash"

	if IsBlockchainValid(bc) {
		t.Fatal("broken previous hash was accepted")
	}
}

func TestTamperedNonce(t *testing.T) {

	bc := CreateBlockchain()

	bc.AddBlock([]Transaction{})

	bc.Blocks[1].Nonce++

	if IsBlockchainValid(bc) {
		t.Fatal("tampered nonce was accepted")
	}
}

func TestProofOfWork(t *testing.T) {

	block := CreateBlock(
		1,
		"previous-hash",
		[]Transaction{},
	)

	if !block.HasValidProofOfWork() {
		t.Fatal("valid proof-of-work was rejected")
	}
}