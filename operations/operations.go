package operations

import (
	"echoLearning/structs"
	"fmt"
	"log"
)

func PerformTransactions(client *structs.Client, transaction structs.Transaction) {
	client.TransChan <- transaction
}

func ProcessTransactions(client structs.Client) error {
	client.Account.Lock()
	defer client.Account.Unlock()

	for transaction := range client.TransChan {

		if transaction.Amount <= 0 {
			return fmt.Errorf("Client %d cannot cashout a non-positive amount. Balance: %d\n", client.ID, client.Account.Balance)
		}

		if transaction.IsDebit {
			if transaction.Amount > client.Account.Balance {
				return fmt.Errorf("Client %d cannot cashout. Insufficient balance. Balance: %d\n", client.ID, client.Account.Balance)
			}

			client.Account.Balance -= transaction.Amount
			log.Printf("Client %d has made cashout of %d. Balance: %d\n", client.ID, transaction.Amount, client.Account.Balance)
		} else {
			client.Account.Balance += transaction.Amount

			log.Printf("Client %d has made deposit of %d. Balance: %d\n", client.ID, transaction.Amount, client.Account.Balance)
		}
	}
	return nil
}

func SendFunds(sender *structs.Client, receiver *structs.Client, amount int) error {
	sender.Account.Lock()
	defer sender.Account.Unlock()

	receiver.Account.Lock()
	defer receiver.Account.Unlock()

	if amount <= 0 {
		return fmt.Errorf("incorrect amount")
	}

	if sender.Account.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	sender.Account.Balance -= amount
	receiver.Account.Balance += amount
	log.Printf("Client %d transferred to client %d %d money\n", sender.ID, receiver.ID, amount)

	return nil
}
