package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

// ============================================================
// MAIN
// ============================================================

func main() {

	command := "status"

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {

	case "status":
		showStatus()

	case "balance":
		showBalance()

	case "wallet":
		showWallet()

	case "wallets":
		showWallets()

	case "createwallet":
		createWallet()

	case "mine":
		mineLBN()

	case "send":
		sendLBN()

	case "help":
		showHelp()

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println()
		showHelp()
	}
}

// ============================================================
// LOAD BLOCKCHAIN
// ============================================================

func loadLBNBlockchain() (Blockchain, error) {

	if BlockchainExists() {
		return LoadBlockchain()
	}

	blockchain := CreateBlockchain()

	if err := SaveBlockchain(blockchain); err != nil {
		return Blockchain{}, err
	}

	return blockchain, nil
}

// ============================================================
// LOAD WALLETS
// ============================================================

func loadWallets() (*Wallet, *Wallet, error) {

	miner, err := LoadOrCreateWallet(minerWalletFile)

	if err != nil {
		return nil, nil, err
	}

	recipient, err := LoadOrCreateWallet(recipientWalletFile)

	if err != nil {
		return nil, nil, err
	}

	writeTextFile(
		"miner_address.txt",
		miner.Address,
	)

	writeTextFile(
		"recipient_address.txt",
		recipient.Address,
	)

	return miner, recipient, nil
}

// ============================================================
// STATUS
// ============================================================
func showStatus() {

	blockchain, err := loadLBNBlockchain()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	utxoSet, err := RebuildUTXOSet(blockchain)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("===== LBN NODE STATUS =====")
	fmt.Println()

	fmt.Println("Blockchain valid:", IsBlockchainValid(blockchain))
	fmt.Println("Blockchain height:", len(blockchain.Blocks)-1)
	fmt.Println("Active UTXOs:", len(utxoSet.Outputs))

	fmt.Println()
	fmt.Println("===== BALANCES =====")
	fmt.Println()

	// Show the balances of all known wallets.
	files, err := os.ReadDir(".")

	if err != nil {
		fmt.Println("ERROR reading wallet files:", err)
		return
	}

	foundWallet := false

	for _, file := range files {

		name := file.Name()

		if len(name) < 12 {
			continue
		}

		if len(name) >= 4 &&
			name[len(name)-4:] == ".key" &&
			len(name) > 12 &&
			name[:7] == "wallet_" {

			wallet, err := LoadWallet(name)

			if err != nil {
				continue
			}

			fmt.Println(
				name+":",
				FormatLBN(
					utxoSet.Balance(wallet.Address),
				),
			)

			fmt.Println(
				"Address:",
				wallet.Address,
			)

			fmt.Println()

			foundWallet = true
		}
	}

	if !foundWallet {
		fmt.Println("No wallets found.")
		fmt.Println()
	}

	currentSupply := CalculateSupply(utxoSet)

	fmt.Println(
		"Current supply:",
		FormatLBN(currentSupply),
	)

	fmt.Println(
		"Maximum supply:",
		FormatLBN(MaxSupplyAtomicUnits()),
	)

	fmt.Println()

	nextHeight := uint64(len(blockchain.Blocks))

	fmt.Println(
		"Next block height:",
		nextHeight,
	)

	fmt.Println(
		"Next block reward:",
		FormatLBN(BlockReward(nextHeight)),
	)

	fmt.Println()

	fmt.Println("Proof of Work: ACTIVE")
	fmt.Println("Mining: READY")
}

// ============================================================
// ALL BALANCES
// ============================================================

func showBalances() {

	blockchain, err := loadLBNBlockchain()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	miner, recipient, err := loadWallets()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	utxoSet, err := RebuildUTXOSet(blockchain)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("===== LBN BALANCES =====")
	fmt.Println()

	fmt.Println(
		"Miner:",
		FormatLBN(
			utxoSet.Balance(miner.Address),
		),
	)

	fmt.Println(
		"Recipient:",
		FormatLBN(
			utxoSet.Balance(recipient.Address),
		),
	)

	fmt.Println()

	fmt.Println(
		"Current supply:",
		FormatLBN(
			CalculateSupply(utxoSet),
		),
	)
}

// ============================================================
// SINGLE WALLET BALANCE
// ============================================================

func showBalance() {

	if len(os.Args) < 3 {

		fmt.Println("Usage:")
		fmt.Println(
			"go run . balance <wallet-address>",
		)

		return
	}

	address := os.Args[2]

	blockchain, err := loadLBNBlockchain()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	utxoSet, err := RebuildUTXOSet(blockchain)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	balance := utxoSet.Balance(address)

	fmt.Println("===== LBN WALLET BALANCE =====")
	fmt.Println()

	fmt.Println(
		"Address:",
		address,
	)

	fmt.Println(
		"Balance:",
		FormatLBN(balance),
	)
}

// ============================================================
// WALLET
// ============================================================

func showWallet() {

	wallet, err := LoadWallet(minerWalletFile)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	wallet.Display()
}

// ============================================================
// CREATE WALLET
// ============================================================

func createWallet() {

	wallet, err := CreateWallet()

	if err != nil {
		fmt.Println(
			"ERROR creating wallet:",
			err,
		)
		return
	}

	filename :=
		"wallet_" +
			wallet.Address[:12] +
			".key"

	if err := SaveWallet(
		wallet,
		filename,
	); err != nil {

		fmt.Println(
			"ERROR saving wallet:",
			err,
		)

		return
	}

	fmt.Println("===== NEW LBN WALLET =====")
	fmt.Println()

	fmt.Println("Address:")
	fmt.Println(wallet.Address)

	fmt.Println()

	fmt.Println("Public Key:")
	fmt.Println(
		hex.EncodeToString(
			wallet.PublicKey,
		),
	)

	fmt.Println()

	fmt.Println("Wallet saved to:")
	fmt.Println(filename)

	fmt.Println()

	fmt.Println("IMPORTANT:")
	fmt.Println(
		"Keep the wallet file secret.",
	)

	fmt.Println(
		"Anyone with the private key can control the wallet.",
	)
}

// ============================================================
// ALL WALLETS
// ============================================================

func showWallets() {

	files, err := os.ReadDir(".")

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	fmt.Println("===== LBN WALLETS =====")
	fmt.Println()

	found := false

	for _, file := range files {

		name := file.Name()

		if len(name) < 12 {
			continue
		}

		if len(name) >= 4 &&
			name[len(name)-4:] == ".key" &&
			len(name) > 12 &&
			name[:7] == "wallet_" {

			wallet, err := LoadWallet(name)

			if err != nil {

				fmt.Println(
					"ERROR loading",
					name,
					":",
					err,
				)

				continue
			}

			fmt.Println(
				"File:",
				name,
			)

			fmt.Println(
				"Address:",
				wallet.Address,
			)

			fmt.Println()

			found = true
		}
	}

	if !found {
		fmt.Println(
			"No additional wallets found.",
		)
	}
}

// ============================================================
// MINE
// ============================================================

func mineLBN() {

	blockchain, err := loadLBNBlockchain()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	miner, _, err := loadWallets()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	nextHeight :=
		uint64(len(blockchain.Blocks))

	reward :=
		BlockReward(nextHeight)

	fmt.Println("===== LBN MINING =====")
	fmt.Println()

	fmt.Println(
		"Mining block:",
		nextHeight,
	)

	fmt.Println(
		"Block reward:",
		FormatLBN(reward),
	)

	fmt.Println()

	block, err := MineBlock(
		&blockchain,
		miner.Address,
		[]Transaction{},
	)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	blockchain.AddBlock(
		block.Transactions,
	)

	if err := SaveBlockchain(
		blockchain,
	); err != nil {

		fmt.Println(
			"ERROR saving blockchain:",
			err,
		)

		return
	}

	utxoSet, err :=
		RebuildUTXOSet(blockchain)

	if err != nil {

		fmt.Println(
			"ERROR rebuilding UTXO set:",
			err,
		)

		return
	}

	if err := SaveUTXOSet(
		utxoSet,
	); err != nil {

		fmt.Println(
			"ERROR saving UTXO set:",
			err,
		)

		return
	}

	fmt.Println(
		"Block mined successfully.",
	)

	fmt.Println()

	fmt.Println(
		"Block:",
		block.Index,
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
		"Proof of Work:",
		block.HasValidProofOfWork(),
	)

	fmt.Println()

	fmt.Println(
		"Miner balance:",
		FormatLBN(
			utxoSet.Balance(
				miner.Address,
			),
		),
	)

	fmt.Println(
		"Current supply:",
		FormatLBN(
			CalculateSupply(utxoSet),
		),
	)
}

// ============================================================
// SEND
// ============================================================

func sendLBN() {

	if len(os.Args) < 5 {

		fmt.Println("Usage:")
		fmt.Println(
			"go run . send <wallet-file> <recipient-address> <amount>",
		)

		return
	}

	senderWalletFile := os.Args[2]
	recipientAddress := os.Args[3]
	amountText := os.Args[4]

	amountLBN, err :=
		strconv.ParseUint(
			amountText,
			10,
			64,
		)

	if err != nil ||
		amountLBN == 0 {

		fmt.Println(
			"ERROR: Invalid amount.",
		)

		return
	}

	if recipientAddress == "" {

		fmt.Println(
			"ERROR: Recipient address is empty.",
		)

		return
	}

	amount :=
		amountLBN *
			1_000_000_000

	sender, err :=
		LoadWallet(senderWalletFile)

	if err != nil {

		fmt.Println(
			"ERROR loading sender wallet:",
			err,
		)

		return
	}

	blockchain, err :=
		loadLBNBlockchain()

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	utxoSet, err :=
		RebuildUTXOSet(blockchain)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	balance :=
		utxoSet.Balance(
			sender.Address,
		)

	fmt.Println("===== LBN SEND =====")
	fmt.Println()

	fmt.Println(
		"From:",
		sender.Address,
	)

	fmt.Println(
		"To:",
		recipientAddress,
	)

	fmt.Println(
		"Amount:",
		FormatLBN(amount),
	)

	fmt.Println(
		"Sender balance:",
		FormatLBN(balance),
	)

	fmt.Println()

	if balance < amount {

		fmt.Println(
			"ERROR: Insufficient balance.",
		)

		return
	}

	tx, err :=
		CreateTransaction(
			sender,
			utxoSet,
			recipientAddress,
			amount,
		)

	if err != nil {

		fmt.Println(
			"ERROR creating transaction:",
			err,
		)

		return
	}

	block, err :=
		MineBlock(
			&blockchain,
			sender.Address,
			[]Transaction{tx},
		)

	if err != nil {

		fmt.Println(
			"ERROR mining transaction:",
			err,
		)

		return
	}

	blockchain.AddBlock(
		block.Transactions,
	)

	if err := SaveBlockchain(
		blockchain,
	); err != nil {

		fmt.Println(
			"ERROR saving blockchain:",
			err,
		)

		return
	}

	utxoSet, err =
		RebuildUTXOSet(blockchain)

	if err != nil {

		fmt.Println(
			"ERROR rebuilding UTXO set:",
			err,
		)

		return
	}

	if err := SaveUTXOSet(
		utxoSet,
	); err != nil {

		fmt.Println(
			"ERROR saving UTXO set:",
			err,
		)

		return
	}

	fmt.Println(
		"Transaction mined successfully.",
	)

	fmt.Println()

	fmt.Println(
		"Block:",
		block.Index,
	)

	fmt.Println(
		"Block hash:",
		block.Hash,
	)

	fmt.Println()

	fmt.Println(
		"Sender balance:",
		FormatLBN(
			utxoSet.Balance(
				sender.Address,
			),
		),
	)

	fmt.Println(
		"Recipient balance:",
		FormatLBN(
			utxoSet.Balance(
				recipientAddress,
			),
		),
	)

	fmt.Println()

	fmt.Println(
		"Current supply:",
		FormatLBN(
			CalculateSupply(utxoSet),
		),
	)
}

// ============================================================
// HELP
// ============================================================

func showHelp() {

	fmt.Println("===== LBN COMMANDS =====")
	fmt.Println()

	fmt.Println(
		"go run . status - Show blockchain status",
	)

	fmt.Println(
		"go run . balance <wallet-address> - Show a wallet balance",
	)

	fmt.Println(
		"go run . wallet - Show miner wallet",
	)

	fmt.Println(
		"go run . wallets - Show created wallets",
	)

	fmt.Println(
		"go run . createwallet - Create a new wallet",
	)

	fmt.Println(
		"go run . mine - Mine a new LBN block",
	)

	fmt.Println(
		"go run . send <wallet-file> <recipient-address> <amount> - Send LBN",
	)

	fmt.Println(
		"go run . help - Show commands",
	)
}

// ============================================================
// TEXT FILE HELPER
// ============================================================

func writeTextFile(
	filename string,
	text string,
) {

	err := os.WriteFile(
		filename,
		[]byte(text+"\n"),
		0600,
	)

	if err != nil {

		fmt.Println(
			"WARNING: Could not write",
			filename,
			":",
			err,
		)
	}
}
