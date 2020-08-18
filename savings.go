package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	bnkdev "github.com/kyleconroy/bnkdev-go"
)

var realMode = flag.Bool("real", false, "create actual transfers")

func main() {
	flag.Parse()

	checking := os.Getenv("BNKDEV_CHECKING_ACCOUNT")
	savings := os.Getenv("BNKDEV_SAVINGS_ACCOUNT")
	cardRoute := os.Getenv("BNKDEV_CARD_ROUTE")
	client := bnkdev.NewClient(os.Getenv("BNKDEV_API_KEY"))
	ctx := context.Background()

	transactions, err := client.ListTransactions(ctx, &bnkdev.ListTransactionsRequest{
		AccountID: savings,
	})
	if err != nil {
		log.Fatal(err)
	}

	seen := map[string]struct{}{}
	for _, tx := range transactions.Data {
		if !strings.HasPrefix("Keep the change: ", tx.Description) {
			continue
		}
		seen[strings.TrimPrefix("Keep the change: ", tx.Description)] = struct{}{}
	}

	transactions, err = client.ListTransactions(ctx, &bnkdev.ListTransactionsRequest{
		AccountID: checking,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range transactions.Data {
		if _, found := seen[tx.ID]; found {
			fmt.Printf("SKIP Already created transfer for %s\n", tx.ID)
			continue
		}
		if tx.RouteID != cardRoute {
			continue
		}
		if tx.Amount >= 0 {
			continue
		}

		change := 100 + (tx.Amount % 100)

		if !*realMode {
			fmt.Printf("SKIP Transferring %d cents to %s\n", change, savings)
			continue
		} else {
			fmt.Printf("Transferring %d cents to %s\n", change, savings)
		}

		// Create an account transfer
		_, err := client.CreateAccountTransfer(ctx, &bnkdev.CreateAccountTransferRequest{
			AccountID:            checking,
			DestinationAccountID: savings,
			Amount:               change,
			Description:          fmt.Sprintf("Keep the change: %s", tx.ID),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
