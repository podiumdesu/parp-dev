package client

import (
	"sync"
)

func (c *Client) WaitForChannelId(wg *sync.WaitGroup, hubSend chan<- []byte) {
	defer wg.Done()

}

// func (c *Client) waitForTransaction(wg *sync.WaitGroup, txHash common.Hash) error {
// 	defer wg.Done()
// 	bcClient, err := c.ConnectToBlockchain()
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		receipt, err := bcClient.TransactionReceipt(context.Background(), txHash)
// 		if err != nil {
// 			if err == ethereum.NotFound {
// 				// Transaction is still pending
// 				time.Sleep(1 * time.Second)
// 				log.Println("Transaction pending:", txHash.Hex())
// 				continue
// 			}
// 			return err // Other errors
// 		}
// 		// Transaction is confirmed
// 		if receipt.Status == 1 {
// 			log.Println("Transaction confirmed:", txHash.Hex())
// 			return nil
// 		}
// 		return fmt.Errorf("transaction failed: %s", txHash.Hex())
// 	}

// }
