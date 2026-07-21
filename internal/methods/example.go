package methods

import (
	"fmt"
	"strings"
)

func Example() string {
	var out strings.Builder

	account := &Account{ID: 1, Owner: "Maria", Balance: 1_000, Active: true, Currency: "RUB"}
	_ = account.Deposit(500)
	_ = account.Withdraw(300)
	copy := account.Renamed("Masha")
	fmt.Fprintf(&out, "account=%s\n", account.Label())
	fmt.Fprintf(&out, "original=%s copy=%s\n", account.Owner, copy.Owner)

	cart := &Cart{MaxItems: 5}
	_ = cart.Add("book", 2)
	_ = cart.Add("pen", 1)
	fmt.Fprintf(&out, "cart=%d", cart.TotalItems())

	return out.String()
}
