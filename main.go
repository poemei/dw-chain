//dw-chain/main.go
package main

import (
    "fmt"
    "os"
    "time"

    "dw-chain/api"
    "dw-chain/blockchain"
)

const dataDir = "data/"

func main() {
    fmt.Println("?? Booting Chain Miner Node...")

    chain := blockchain.InitBlockchain()
    config := blockchain.AppConfig
    balances := blockchain.LoadBalances()
    blockchain.LoadPeers()

    // Start local status server
    api.StartServer(&chain, &balances)

    for {
        fmt.Println("?? Checking transactions...")

        txs := blockchain.LoadTransactions()
        if len(txs) == 0 {
            fmt.Println("??? No pending transactions. Sleeping...")
            time.Sleep(15 * time.Second)
            continue
        }

        validTxs := []blockchain.Transaction{}
        for _, tx := range txs {
            if tx.Sender != "GENESIS" && balances[tx.Sender] < tx.Amount {
                fmt.Printf("? Tx skipped (insufficient funds): %s", tx.Sender)
                continue
            }
            balances[tx.Sender] -= tx.Amount
            balances[tx.Recipient] += tx.Amount
            validTxs = append(validTxs, tx)
        }

        if len(validTxs) == 0 {
            fmt.Println("?? No valid transactions. Sleeping...")
            time.Sleep(15 * time.Second)
            continue
        }

        newDifficulty := chain.AdjustDifficulty()
        newBlock := blockchain.Block{
            Index:        len(chain.Blocks),
            Timestamp:    blockchain.Now(),
            Transactions: validTxs,
            PrevHash:     chain.LatestBlock().Hash,
            Difficulty:   newDifficulty,
            Miner:        "PiNode",
        }

        fmt.Println("?? Starting mining process...")
        blockchain.MineBlock(&newBlock, config.MinerThreads)
        fmt.Printf("? Block mined: %s", newBlock.Hash)

        chain.Blocks = append(chain.Blocks, newBlock)
        chain.Save()

        blockchain.SaveBalances(balances)
        os.WriteFile(dataDir+"transactions.json", []byte("[]"), 0644)

        fmt.Println("?? Block committed. Pausing...")
        blockchain.BroadcastBlock(newBlock)
        time.Sleep(15 * time.Second)
    }
}