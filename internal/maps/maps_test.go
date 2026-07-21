package maps

import (
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  map[string]int
	}{
		{"nil", nil, map[string]int{}},
		{"empty", []string{}, map[string]int{}},
		{"one", []string{"go"}, map[string]int{"go": 1}},
		{"duplicates", []string{"go", "map", "go"}, map[string]int{"go": 2, "map": 1}},
		{"case-sensitive", []string{"Go", "go", "GO"}, map[string]int{"Go": 1, "go": 1, "GO": 1}},
		{"empty-word", []string{"", ""}, map[string]int{"": 2}},
		{"spaces", []string{"go lang", "go lang"}, map[string]int{"go lang": 2}},
		{"unicode", []string{"го", "go", "го"}, map[string]int{"го": 2, "go": 1}},
		{"many", []string{"a", "b", "c", "a", "b", "a"}, map[string]int{"a": 3, "b": 2, "c": 1}},
		{"punctuation", []string{"go!", "go", "go!"}, map[string]int{"go!": 2, "go": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountWords(tt.words)
			if got == nil {
				t.Fatal("result map must be initialized")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		values   map[string]int
		key      string
		fallback int
		want     int
	}{
		{"nil", nil, "a", 7, 7},
		{"empty", map[string]int{}, "a", 7, 7},
		{"existing-positive", map[string]int{"a": 3}, "a", 7, 3},
		{"existing-zero", map[string]int{"a": 0}, "a", 7, 0},
		{"existing-negative", map[string]int{"a": -2}, "a", 7, -2},
		{"missing-zero-fallback", map[string]int{"a": 1}, "b", 0, 0},
		{"empty-key-existing", map[string]int{"": 9}, "", 1, 9},
		{"empty-key-missing", map[string]int{"a": 9}, "", 1, 1},
		{"unicode", map[string]int{"ключ": 5}, "ключ", 1, 5},
		{"case-sensitive", map[string]int{"Go": 5}, "go", 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOrDefault(tt.values, tt.key, tt.fallback); got != tt.want {
				t.Fatalf("got %d want %d", got, tt.want)
			}
		})
	}
}

func TestDeleteAndReport(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int
		key    string
		want   bool
		left   map[string]int
	}{
		{"nil", nil, "a", false, nil},
		{"empty", map[string]int{}, "a", false, map[string]int{}},
		{"existing", map[string]int{"a": 1}, "a", true, map[string]int{}},
		{"missing", map[string]int{"a": 1}, "b", false, map[string]int{"a": 1}},
		{"existing-zero", map[string]int{"a": 0}, "a", true, map[string]int{}},
		{"empty-key", map[string]int{"": 1}, "", true, map[string]int{}},
		{"case-sensitive", map[string]int{"Go": 1}, "go", false, map[string]int{"Go": 1}},
		{"keep-other", map[string]int{"a": 1, "b": 2}, "a", true, map[string]int{"b": 2}},
		{"unicode", map[string]int{"ключ": 3}, "ключ", true, map[string]int{}},
		{"negative-value", map[string]int{"a": -1}, "a", true, map[string]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteAndReport(tt.values, tt.key); got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if !reflect.DeepEqual(tt.values, tt.left) {
				t.Fatalf("map after call %#v want %#v", tt.values, tt.left)
			}
		})
	}
}

func TestMergeCounters(t *testing.T) {
	tests := []struct {
		name        string
		left, right map[string]int
		want        map[string]int
	}{
		{"both-nil", nil, nil, map[string]int{}},
		{"left-only", map[string]int{"a": 1}, nil, map[string]int{"a": 1}},
		{"right-only", nil, map[string]int{"b": 2}, map[string]int{"b": 2}},
		{"disjoint", map[string]int{"a": 1}, map[string]int{"b": 2}, map[string]int{"a": 1, "b": 2}},
		{"overlap", map[string]int{"a": 1}, map[string]int{"a": 2}, map[string]int{"a": 3}},
		{"zeros", map[string]int{"a": 0}, map[string]int{"a": 0}, map[string]int{"a": 0}},
		{"negative", map[string]int{"a": 5}, map[string]int{"a": -2}, map[string]int{"a": 3}},
		{"empty-key", map[string]int{"": 1}, map[string]int{"": 4}, map[string]int{"": 5}},
		{"many", map[string]int{"a": 1, "b": 2}, map[string]int{"b": 3, "c": 4}, map[string]int{"a": 1, "b": 5, "c": 4}},
		{"case-sensitive", map[string]int{"Go": 1}, map[string]int{"go": 2}, map[string]int{"Go": 1, "go": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftBefore := cloneStringIntMap(tt.left)
			rightBefore := cloneStringIntMap(tt.right)
			got := MergeCounters(tt.left, tt.right)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.left, leftBefore) || !reflect.DeepEqual(tt.right, rightBefore) {
				t.Fatal("input maps were changed")
			}
		})
	}
}

func TestKeysByValueAtLeast(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int
		min    int
		want   []string
	}{
		{"nil", nil, 0, []string{}},
		{"empty", map[string]int{}, 0, []string{}},
		{"all", map[string]int{"b": 2, "a": 1}, 1, []string{"a", "b"}},
		{"none", map[string]int{"a": 1}, 2, []string{}},
		{"boundary", map[string]int{"a": 2, "b": 1}, 2, []string{"a"}},
		{"negative-min", map[string]int{"a": -2, "b": 0}, -1, []string{"b"}},
		{"negative-values", map[string]int{"a": -2, "b": -1}, -2, []string{"a", "b"}},
		{"empty-key", map[string]int{"": 5, "a": 1}, 5, []string{""}},
		{"case-order", map[string]int{"go": 2, "Go": 2}, 2, []string{"Go", "go"}},
		{"unicode", map[string]int{"я": 3, "а": 3}, 3, []string{"а", "я"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeysByValueAtLeast(tt.values, tt.min); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestMaxKeyByValue(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int
		key    string
		ok     bool
	}{
		{"nil", nil, "", false},
		{"empty", map[string]int{}, "", false},
		{"one", map[string]int{"a": 1}, "a", true},
		{"positive", map[string]int{"a": 1, "b": 3}, "b", true},
		{"tie", map[string]int{"b": 3, "a": 3}, "a", true},
		{"all-negative", map[string]int{"a": -3, "b": -1}, "b", true},
		{"zero-wins", map[string]int{"a": -1, "b": 0}, "b", true},
		{"empty-key", map[string]int{"": 4, "a": 3}, "", true},
		{"case-tie", map[string]int{"go": 2, "Go": 2}, "Go", true},
		{"unicode", map[string]int{"я": 5, "а": 5}, "а", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, ok := MaxKeyByValue(tt.values)
			if key != tt.key || ok != tt.ok {
				t.Fatalf("got %q,%t want %q,%t", key, ok, tt.key, tt.ok)
			}
		})
	}
}

func TestBuildUserIndex(t *testing.T) {
	tests := []struct {
		name  string
		users []User
		want  map[int]User
	}{
		{"nil", nil, map[int]User{}},
		{"empty", []User{}, map[int]User{}},
		{"one", []User{{ID: 1, Name: "A"}}, map[int]User{1: {ID: 1, Name: "A"}}},
		{"two", []User{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, map[int]User{1: {ID: 1, Name: "A"}, 2: {ID: 2, Name: "B"}}},
		{"duplicate", []User{{ID: 1, Name: "A"}, {ID: 1, Name: "B"}}, map[int]User{1: {ID: 1, Name: "B"}}},
		{"zero-id", []User{{ID: 0, Name: "Zero"}}, map[int]User{0: {ID: 0, Name: "Zero"}}},
		{"negative-id", []User{{ID: -1, Name: "N"}}, map[int]User{-1: {ID: -1, Name: "N"}}},
		{"keep-fields", []User{{ID: 1, City: "M", Active: true}}, map[int]User{1: {ID: 1, City: "M", Active: true}}},
		{"tags", []User{{ID: 1, Tags: []string{"go"}}}, map[int]User{1: {ID: 1, Tags: []string{"go"}}}},
		{"last-full", []User{{ID: 1, Name: "A", Active: true}, {ID: 1, Name: "B", Active: false}}, map[int]User{1: {ID: 1, Name: "B", Active: false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildUserIndex(tt.users)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestGroupActiveUsersByCity(t *testing.T) {
	tests := []struct {
		name  string
		users []User
		want  map[string][]User
	}{
		{"nil", nil, map[string][]User{}},
		{"empty", []User{}, map[string][]User{}},
		{"inactive-only", []User{{Name: "A", City: "M"}}, map[string][]User{}},
		{"one-active", []User{{Name: "A", City: "M", Active: true}}, map[string][]User{"M": {{Name: "A", City: "M", Active: true}}}},
		{"same-city", []User{{Name: "A", City: "M", Active: true}, {Name: "B", City: "M", Active: true}}, map[string][]User{"M": {{Name: "A", City: "M", Active: true}, {Name: "B", City: "M", Active: true}}}},
		{"preserve-order", []User{{Name: "B", City: "M", Active: true}, {Name: "A", City: "M", Active: true}}, map[string][]User{"M": {{Name: "B", City: "M", Active: true}, {Name: "A", City: "M", Active: true}}}},
		{"two-cities", []User{{Name: "A", City: "M", Active: true}, {Name: "B", City: "B", Active: true}}, map[string][]User{"M": {{Name: "A", City: "M", Active: true}}, "B": {{Name: "B", City: "B", Active: true}}}},
		{"mixed", []User{{Name: "A", City: "M", Active: false}, {Name: "B", City: "M", Active: true}}, map[string][]User{"M": {{Name: "B", City: "M", Active: true}}}},
		{"empty-city", []User{{Name: "A", City: "", Active: true}}, map[string][]User{"": {{Name: "A", City: "", Active: true}}}},
		{"keep-tags", []User{{Name: "A", City: "M", Active: true, Tags: []string{"go"}}}, map[string][]User{"M": {{Name: "A", City: "M", Active: true, Tags: []string{"go"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupActiveUsersByCity(tt.users)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestCountTags(t *testing.T) {
	tests := []struct {
		name  string
		users []User
		want  map[string]int
	}{
		{"nil", nil, map[string]int{}},
		{"empty", []User{}, map[string]int{}},
		{"no-tags", []User{{Name: "A"}}, map[string]int{}},
		{"one", []User{{Tags: []string{"go"}}}, map[string]int{"go": 1}},
		{"duplicate-one-user", []User{{Tags: []string{"go", "go"}}}, map[string]int{"go": 2}},
		{"duplicate-users", []User{{Tags: []string{"go"}}, {Tags: []string{"go"}}}, map[string]int{"go": 2}},
		{"many", []User{{Tags: []string{"go", "api"}}, {Tags: []string{"sql", "go"}}}, map[string]int{"go": 2, "api": 1, "sql": 1}},
		{"empty-tag", []User{{Tags: []string{""}}}, map[string]int{"": 1}},
		{"case-sensitive", []User{{Tags: []string{"Go", "go"}}}, map[string]int{"Go": 1, "go": 1}},
		{"unicode", []User{{Tags: []string{"бэк", "бэк"}}}, map[string]int{"бэк": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountTags(tt.users)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestBuildInventory(t *testing.T) {
	tests := []struct {
		name     string
		products []Product
		want     map[string]Product
	}{
		{"nil", nil, map[string]Product{}},
		{"empty", []Product{}, map[string]Product{}},
		{"one", []Product{{SKU: "a", Quantity: 1, Price: 10}}, map[string]Product{"a": {SKU: "a", Quantity: 1, Price: 10}}},
		{"two", []Product{{SKU: "a"}, {SKU: "b"}}, map[string]Product{"a": {SKU: "a"}, "b": {SKU: "b"}}},
		{"duplicate", []Product{{SKU: "a", Price: 1}, {SKU: "a", Price: 2}}, map[string]Product{"a": {SKU: "a", Price: 2}}},
		{"empty-sku", []Product{{SKU: "", Price: 1}}, map[string]Product{"": {SKU: "", Price: 1}}},
		{"negative-quantity", []Product{{SKU: "a", Quantity: -1}}, map[string]Product{"a": {SKU: "a", Quantity: -1}}},
		{"zero-price", []Product{{SKU: "a", Price: 0}}, map[string]Product{"a": {SKU: "a", Price: 0}}},
		{"case-sensitive", []Product{{SKU: "A"}, {SKU: "a"}}, map[string]Product{"A": {SKU: "A"}, "a": {SKU: "a"}}},
		{"last-full", []Product{{SKU: "a", Quantity: 1, Price: 2}, {SKU: "a", Quantity: 3, Price: 4}}, map[string]Product{"a": {SKU: "a", Quantity: 3, Price: 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildInventory(tt.products)
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestReserveStock(t *testing.T) {
	tests := []struct {
		name      string
		inventory map[string]Product
		sku       string
		count     int
		want      bool
		quantity  int
	}{
		{"nil", nil, "a", 1, false, 0},
		{"missing", map[string]Product{}, "a", 1, false, 0},
		{"zero-count", map[string]Product{"a": {SKU: "a", Quantity: 5}}, "a", 0, false, 5},
		{"negative-count", map[string]Product{"a": {SKU: "a", Quantity: 5}}, "a", -1, false, 5},
		{"not-enough", map[string]Product{"a": {SKU: "a", Quantity: 2}}, "a", 3, false, 2},
		{"exact", map[string]Product{"a": {SKU: "a", Quantity: 2}}, "a", 2, true, 0},
		{"partial", map[string]Product{"a": {SKU: "a", Quantity: 5}}, "a", 2, true, 3},
		{"zero-stock", map[string]Product{"a": {SKU: "a", Quantity: 0}}, "a", 1, false, 0},
		{"keep-price", map[string]Product{"a": {SKU: "a", Quantity: 5, Price: 9}}, "a", 1, true, 4},
		{"case-sensitive", map[string]Product{"A": {SKU: "A", Quantity: 5}}, "a", 1, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReserveStock(tt.inventory, tt.sku, tt.count)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if product, ok := tt.inventory[tt.sku]; ok && product.Quantity != tt.quantity {
				t.Fatalf("quantity=%d want %d", product.Quantity, tt.quantity)
			}
		})
	}
}

func TestRestock(t *testing.T) {
	product := func(quantity int) *Product { return &Product{SKU: "a", Quantity: quantity, Price: 10} }
	tests := []struct {
		name      string
		inventory map[string]*Product
		sku       string
		count     int
		want      bool
		quantity  int
	}{
		{"nil-map", nil, "a", 1, false, 0},
		{"missing", map[string]*Product{}, "a", 1, false, 0},
		{"nil-product", map[string]*Product{"a": nil}, "a", 1, false, 0},
		{"zero-count", map[string]*Product{"a": product(5)}, "a", 0, false, 5},
		{"negative-count", map[string]*Product{"a": product(5)}, "a", -1, false, 5},
		{"add-one", map[string]*Product{"a": product(5)}, "a", 1, true, 6},
		{"add-many", map[string]*Product{"a": product(0)}, "a", 10, true, 10},
		{"negative-stock", map[string]*Product{"a": product(-2)}, "a", 3, true, 1},
		{"case-sensitive", map[string]*Product{"A": product(2)}, "a", 1, false, 0},
		{"keep-pointer", map[string]*Product{"a": product(2)}, "a", 2, true, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := tt.inventory[tt.sku]
			got := Restock(tt.inventory, tt.sku, tt.count)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if before != nil && got {
				if tt.inventory[tt.sku] != before {
					t.Fatal("pointer must stay the same")
				}
				if before.Quantity != tt.quantity {
					t.Fatalf("quantity=%d want %d", before.Quantity, tt.quantity)
				}
			}
		})
	}
}

func TestLowStockSKUs(t *testing.T) {
	tests := []struct {
		name      string
		inventory map[string]Product
		limit     int
		want      []string
	}{
		{"nil", nil, 1, []string{}},
		{"empty", map[string]Product{}, 1, []string{}},
		{"none", map[string]Product{"a": {Quantity: 2}}, 1, []string{}},
		{"one", map[string]Product{"a": {Quantity: 1}}, 1, []string{"a"}},
		{"boundary", map[string]Product{"a": {Quantity: 2}, "b": {Quantity: 1}}, 2, []string{"a", "b"}},
		{"sorted", map[string]Product{"b": {Quantity: 0}, "a": {Quantity: 0}}, 0, []string{"a", "b"}},
		{"negative-limit", map[string]Product{"a": {Quantity: -1}, "b": {Quantity: 0}}, -1, []string{"a"}},
		{"negative-stock", map[string]Product{"a": {Quantity: -2}}, 0, []string{"a"}},
		{"empty-sku", map[string]Product{"": {Quantity: 0}}, 0, []string{""}},
		{"case-order", map[string]Product{"a": {Quantity: 0}, "A": {Quantity: 0}}, 0, []string{"A", "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LowStockSKUs(tt.inventory, tt.limit); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestInventoryValue(t *testing.T) {
	tests := []struct {
		name      string
		inventory map[string]Product
		want      int
	}{
		{"nil", nil, 0},
		{"empty", map[string]Product{}, 0},
		{"one", map[string]Product{"a": {Quantity: 2, Price: 10}}, 20},
		{"two", map[string]Product{"a": {Quantity: 2, Price: 10}, "b": {Quantity: 3, Price: 5}}, 35},
		{"zero-quantity", map[string]Product{"a": {Quantity: 0, Price: 10}}, 0},
		{"zero-price", map[string]Product{"a": {Quantity: 2, Price: 0}}, 0},
		{"negative-quantity", map[string]Product{"a": {Quantity: -2, Price: 10}}, -20},
		{"negative-price", map[string]Product{"a": {Quantity: 2, Price: -10}}, -20},
		{"mixed", map[string]Product{"a": {Quantity: 2, Price: 10}, "b": {Quantity: -1, Price: 5}}, 15},
		{"large", map[string]Product{"a": {Quantity: 1000, Price: 2000}}, 2_000_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InventoryValue(tt.inventory); got != tt.want {
				t.Fatalf("got %d want %d", got, tt.want)
			}
		})
	}
}

func TestApplyPriceUpdates(t *testing.T) {
	tests := []struct {
		name      string
		inventory map[string]Product
		updates   map[string]int
		wantCount int
		want      map[string]Product
	}{
		{"nil-both", nil, nil, 0, nil},
		{"nil-inventory", nil, map[string]int{"a": 10}, 0, nil},
		{"nil-updates", map[string]Product{"a": {SKU: "a", Price: 1}}, nil, 0, map[string]Product{"a": {SKU: "a", Price: 1}}},
		{"existing", map[string]Product{"a": {SKU: "a", Price: 1}}, map[string]int{"a": 10}, 1, map[string]Product{"a": {SKU: "a", Price: 10}}},
		{"missing", map[string]Product{"a": {SKU: "a", Price: 1}}, map[string]int{"b": 10}, 0, map[string]Product{"a": {SKU: "a", Price: 1}}},
		{"zero-price", map[string]Product{"a": {SKU: "a", Price: 1}}, map[string]int{"a": 0}, 0, map[string]Product{"a": {SKU: "a", Price: 1}}},
		{"negative-price", map[string]Product{"a": {SKU: "a", Price: 1}}, map[string]int{"a": -1}, 0, map[string]Product{"a": {SKU: "a", Price: 1}}},
		{"many", map[string]Product{"a": {SKU: "a", Price: 1}, "b": {SKU: "b", Price: 2}}, map[string]int{"a": 10, "b": 20}, 2, map[string]Product{"a": {SKU: "a", Price: 10}, "b": {SKU: "b", Price: 20}}},
		{"mixed", map[string]Product{"a": {SKU: "a", Price: 1}, "b": {SKU: "b", Price: 2}}, map[string]int{"a": 10, "b": 0, "c": 30}, 1, map[string]Product{"a": {SKU: "a", Price: 10}, "b": {SKU: "b", Price: 2}}},
		{"keep-fields", map[string]Product{"a": {SKU: "a", Quantity: 3, Price: 1}}, map[string]int{"a": 10}, 1, map[string]Product{"a": {SKU: "a", Quantity: 3, Price: 10}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplyPriceUpdates(tt.inventory, tt.updates); got != tt.wantCount {
				t.Fatalf("count=%d want %d", got, tt.wantCount)
			}
			if !reflect.DeepEqual(tt.inventory, tt.want) {
				t.Fatalf("inventory=%#v want %#v", tt.inventory, tt.want)
			}
		})
	}
}

func cloneStringIntMap(src map[string]int) map[string]int {
	if src == nil {
		return nil
	}
	result := make(map[string]int, len(src))
	for key, value := range src {
		result[key] = value
	}
	return result
}
