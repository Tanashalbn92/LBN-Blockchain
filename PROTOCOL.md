# LBN — Lobban Network Protocol

Version: 0.1.0

## 1. Purpose

LBN is an independent cryptocurrency intended to operate on its own
peer-to-peer blockchain network.

LBN is not an ERC-20 token or a token hosted on another blockchain.

The LBN protocol will define its own:

- Blockchain
- Transactions
- Wallets
- Monetary supply
- Consensus mechanism
- Peer-to-peer network
- Block validation rules

---

## 2. Native Currency

Name: LBN

Smallest unit: nanoLBN

1 LBN = 1,000,000,000 nanoLBN

All amounts inside the protocol will be represented internally as
integer nanoLBN values.

Floating-point numbers must never be used for monetary calculations.

---

## 3. Maximum Supply

Maximum supply:

1,000,000,000 LBN

The protocol must never create more than the maximum supply.

The maximum supply is a consensus rule.

---

## 4. Consensus

LBN will initially use Proof of Work.

The purpose of Proof of Work is to allow independent participants
to compete to produce valid blocks and secure the network.

No central server will control block production.

---

## 5. Block Time

Target block interval:

2 minutes

The network difficulty will eventually adjust so that blocks are
produced approximately every two minutes over time.

---

## 6. Transactions

LBN will use a UTXO transaction model.

A transaction consumes previously created unspent transaction outputs
and creates new transaction outputs.

Every transaction must satisfy:

total inputs >= total outputs

The difference between inputs and outputs may become a transaction fee.

---

## 7. Transaction Signatures

Transactions must be cryptographically signed by the owner of the
spending key.

Nodes must verify transaction signatures before accepting transactions.

A transaction with an invalid signature must be rejected.

---

## 8. Double Spending

A UTXO may only be spent once.

Nodes must reject transactions attempting to spend an already-spent UTXO.

---

## 9. Coinbase / Block Reward

New LBN may enter circulation through valid block rewards.

Only the consensus rules may create new LBN.

Ordinary transactions cannot create new LBN.

The total amount of newly created LBN must never cause the total supply
to exceed the maximum supply.

---

## 10. Block Reward

Initial block reward:

50 LBN

The block reward will decrease according to the LBN issuance schedule.

The exact reduction schedule will be implemented as a consensus rule
before mainnet launch.

---

## 11. Transaction Fees

Transactions may include fees.

Fees are calculated as:

total inputs - total outputs

Valid block producers may collect valid transaction fees.

---

## 12. Genesis Block

The first block in the LBN blockchain is the Genesis Block.

The Genesis Block has:

Block height: 0

Previous hash:

none

The Genesis Block must be permanently defined before mainnet launch.

---

## 13. Block Validation

Nodes must verify:

- Block index
- Previous block hash
- Block hash
- Timestamp rules
- Proof of Work
- Transaction validity
- Transaction signatures
- UTXO availability
- Block reward
- Transaction fees
- Maximum supply rules

An invalid block must not be accepted into the canonical chain.

---

## 14. Chain Selection

Nodes will follow the valid chain with the greatest accumulated
Proof of Work.

A chain containing invalid blocks must be rejected regardless of
its length.

---

## 15. Network

LBN will operate as a peer-to-peer network.

Nodes will communicate directly with other LBN nodes.

No single node will be required for the network to function.

---

## 16. Mainnet

The LBN mainnet will not launch until:

- Consensus rules are implemented
- Wallet security is implemented
- Transaction validation is implemented
- UTXO accounting is implemented
- Proof of Work is implemented
- Peer-to-peer networking is implemented
- Chain synchronization is implemented
- Double-spend protection is tested
- Supply rules are tested
- Security testing is completed

---

## 17. Protocol Changes

Changes to consensus rules must be versioned.

Consensus changes must never be introduced casually after mainnet launch.

Backward compatibility and network upgrade procedures must be defined
before a production release.

---

## 18. Development Status

Current status:

Protocol design and development.

Mainnet status:

NOT LAUNCHED.

No LBN currently has an established market value.

No LBN should be represented as having a guaranteed monetary value.
