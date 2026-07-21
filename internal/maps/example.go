package maps

import (
	"fmt"
	"strings"
)

func Example() string {
	var out strings.Builder

	words := []string{"go", "map", "go", "test"}
	frequency := CountWords(words)
	fmt.Fprintf(&out, "go=%d missing=%d\n", frequency["go"], GetOrDefault(frequency, "sql", 42))

	users := []User{
		{ID: 1, Name: "Maria", City: "Moscow", Active: true, Tags: []string{"go", "api"}},
		{ID: 2, Name: "Alex", City: "Berlin", Active: false, Tags: []string{"sql"}},
		{ID: 3, Name: "Ada", City: "Moscow", Active: true, Tags: []string{"go"}},
	}
	groups := GroupActiveUsersByCity(users)
	tags := CountTags(users)
	fmt.Fprintf(&out, "moscow=%d tags=go:%d api:%d sql:%d\n", len(groups["Moscow"]), tags["go"], tags["api"], tags["sql"])

	inventory := BuildInventory([]Product{{SKU: "book", Quantity: 3, Price: 100}, {SKU: "pen", Quantity: 10, Price: 20}})
	_ = ReserveStock(inventory, "book", 2)
	fmt.Fprintf(&out, "low=%v value=%d", LowStockSKUs(inventory, 1), InventoryValue(inventory))

	return out.String()
}
