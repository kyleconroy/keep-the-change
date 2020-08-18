# Keep the Change

Inspired by [Bank of America's Keep the Change](https://www.bankofamerica.com/deposits/keep-the-change/) program, this code will round up your purchases to the nearest dollar and transfer the difference to your savings account.

## Configuration

| Environment Variable | Description |
| --- | --- |
| BNKDEV_API_KEY | bnk.dev API Key |
| BNKDEV_CARD_ROUTE | Debit card to use |
| BNKDEV_CHECKING_ACCOUNT | Source account for transactions |
| BNKDEV_SAVINGS_ACCOUNT | Destination account for savings |

## Running

By default, no transfers are created. You can see which actions will be taken by running the following:

```
go run savings.go
```

Once you're happy with the result, pass the `-real` flag to make the transfers.

```
go run savings.go -real
```

Don't worry about duplicate transfers, the code is idempotent and makes sure at most one transfer exists for each transaction.
