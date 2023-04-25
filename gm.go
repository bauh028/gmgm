// SPDX-License-Identifier: MIT
package main

import (
    "context"
    "fmt"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func main() {
    // Connect to the Ethereum network
    client, err := ethclient.Dial("https://mainnet.infura.io/v3/your-infura-api-key")
    if err != nil {
        panic(err)
    }

    // Set up the GM contract instance
    contractAddress := common.HexToAddress("0xYOUR_CONTRACT_ADDRESS")
    GMContract, err := abi.NewGM(contractAddress, client)
    if err != nil {
        panic(err)
    }

    // Set up the sender account
    privateKey := "0xYOUR_PRIVATE_KEY"
    privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
    if err != nil {
        panic(err)
    }
    publicKey := privateKeyECDSA.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

    // Set up the transaction options
    auth := bind.NewKeyedTransactor(privateKeyECDSA)
    auth.Nonce = nil // only on next live network deployment
    auth.Value = big.NewInt(0) // in wei
    auth.GasLimit = uint64(500000) // in units
    auth.GasPrice = big.NewInt(5000000000) // 5 gwei

    // Call the GM contract method to send the message
    message := "Hello GM!"
    tx, err := GMContract.SendMessage(auth, message)
    if err != nil {
        panic(err)
    }

    // Wait for the transaction to be mined
    ctx := context.Background()
    receipt, err := bind.WaitMined(ctx, client, tx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Transaction hash: %v\n", tx.Hash())
    fmt.Printf("Gas used: %v\n", receipt.GasUsed)
}
