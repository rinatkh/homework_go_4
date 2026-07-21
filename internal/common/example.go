package common

import (
	"fmt"
	"strings"

	"github.com/rinatkh/homework_go_4/internal/methods"
)

func Example() string {
	var out strings.Builder

	accounts := map[int]*methods.Account{
		1: {ID: 1, Owner: "Maria", Balance: 1_000, Active: true, Currency: "RUB"},
		2: {ID: 2, Owner: "Alex", Balance: 100, Active: false, Currency: "RUB"},
	}
	result := ApplyOperations(accounts, []Operation{
		{AccountID: 1, Kind: "deposit", Amount: 200},
		{AccountID: 2, Kind: "withdraw", Amount: 50},
	})
	fmt.Fprintf(&out, "balances=%v errors=%d\n", result.Balances, len(result.Errors))

	yes := true
	users, active := ApplyUserEvents(nil, nil, []Event{{UserID: 1, Name: "Maria", Active: &yes}})
	fmt.Fprintf(&out, "users=%v active=%v", users, active)

	return out.String()
}
