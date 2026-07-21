package methods

import (
	"errors"
	"reflect"
	"testing"
)

func TestAccountLabel(t *testing.T) {
	tests := []struct {
		name string
		a    Account
		want string
	}{
		{"zero", Account{}, "#0 : 0 "},
		{"basic", Account{ID: 1, Owner: "Maria", Balance: 100, Currency: "RUB"}, "#1 Maria: 100 RUB"},
		{"negative-id", Account{ID: -1, Owner: "A", Balance: 5, Currency: "USD"}, "#-1 A: 5 USD"},
		{"negative-balance", Account{ID: 2, Owner: "A", Balance: -5, Currency: "EUR"}, "#2 A: -5 EUR"},
		{"empty-owner", Account{ID: 3, Balance: 7, Currency: "RUB"}, "#3 : 7 RUB"},
		{"spaces-owner", Account{ID: 4, Owner: "A B", Balance: 8, Currency: "RUB"}, "#4 A B: 8 RUB"},
		{"unicode", Account{ID: 5, Owner: "Мария", Balance: 9, Currency: "₽"}, "#5 Мария: 9 ₽"},
		{"empty-currency", Account{ID: 6, Owner: "A", Balance: 10}, "#6 A: 10 "},
		{"active-ignored", Account{ID: 7, Owner: "A", Balance: 11, Active: true, Currency: "RUB"}, "#7 A: 11 RUB"},
		{"large", Account{ID: 999, Owner: "Owner", Balance: 1_000_000, Currency: "RUB"}, "#999 Owner: 1000000 RUB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Label(); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestAccountCanWithdraw(t *testing.T) {
	tests := []struct {
		name   string
		a      Account
		amount int
		want   bool
	}{
		{"inactive", Account{Balance: 100}, 10, false},
		{"zero", Account{Balance: 100, Active: true}, 0, false},
		{"negative", Account{Balance: 100, Active: true}, -1, false},
		{"less", Account{Balance: 100, Active: true}, 99, true},
		{"equal", Account{Balance: 100, Active: true}, 100, true},
		{"more", Account{Balance: 100, Active: true}, 101, false},
		{"zero-balance", Account{Balance: 0, Active: true}, 1, false},
		{"negative-balance", Account{Balance: -1, Active: true}, 1, false},
		{"one", Account{Balance: 1, Active: true}, 1, true},
		{"large", Account{Balance: 1_000_000, Active: true}, 999_999, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.CanWithdraw(tt.amount); got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
		})
	}
}

func TestAccountRenamed(t *testing.T) {
	names := []string{"", "A", "B", "Maria", "Мария", "A B", "  ", "Go", "owner-1", "very long owner name"}
	for i, name := range names {
		t.Run(nameForTest(name, i), func(t *testing.T) {
			original := Account{ID: i + 1, Owner: "old", Balance: 100, Active: true, Currency: "RUB"}
			got := original.Renamed(name)
			if original.Owner != "old" {
				t.Fatalf("original changed: %#v", original)
			}
			want := original
			want.Owner = name
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v want %#v", got, want)
			}
		})
	}
}

func TestAccountRename(t *testing.T) {
	tests := []struct {
		name    string
		a       *Account
		owner   string
		wantErr error
	}{
		{"nil", nil, "A", ErrNilAccount},
		{"empty", &Account{Owner: "old"}, "", nil},
		{"basic", &Account{Owner: "old"}, "new", nil},
		{"same", &Account{Owner: "same"}, "same", nil},
		{"unicode", &Account{Owner: "old"}, "Мария", nil},
		{"spaces", &Account{Owner: "old"}, "A B", nil},
		{"only-spaces", &Account{Owner: "old"}, "  ", nil},
		{"keep-balance", &Account{Owner: "old", Balance: 10}, "new", nil},
		{"keep-active", &Account{Owner: "old", Active: true}, "new", nil},
		{"keep-currency", &Account{Owner: "old", Currency: "RUB"}, "new", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before Account
			if tt.a != nil {
				before = *tt.a
			}
			err := tt.a.Rename(tt.owner)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.a != nil {
				if tt.a.Owner != tt.owner {
					t.Fatalf("owner=%q want %q", tt.a.Owner, tt.owner)
				}
				if tt.a.ID != before.ID || tt.a.Balance != before.Balance || tt.a.Active != before.Active || tt.a.Currency != before.Currency {
					t.Fatalf("unrelated fields changed: before=%#v after=%#v", before, *tt.a)
				}
			}
		})
	}
}

func TestAccountDeposit(t *testing.T) {
	tests := []struct {
		name        string
		a           *Account
		amount      int
		wantErr     error
		wantBalance int
	}{
		{"nil", nil, 10, ErrNilAccount, 0},
		{"inactive", &Account{Balance: 100}, 10, ErrInactive, 100},
		{"inactive-invalid", &Account{Balance: 100}, 0, ErrInactive, 100},
		{"zero", &Account{Balance: 100, Active: true}, 0, ErrInvalidAmount, 100},
		{"negative", &Account{Balance: 100, Active: true}, -1, ErrInvalidAmount, 100},
		{"one", &Account{Balance: 0, Active: true}, 1, nil, 1},
		{"positive", &Account{Balance: 100, Active: true}, 50, nil, 150},
		{"negative-balance", &Account{Balance: -10, Active: true}, 5, nil, -5},
		{"large", &Account{Balance: 1_000_000, Active: true}, 1_000_000, nil, 2_000_000},
		{"keep-fields", &Account{ID: 1, Owner: "A", Balance: 10, Active: true, Currency: "RUB"}, 5, nil, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before Account
			if tt.a != nil {
				before = *tt.a
			}
			err := tt.a.Deposit(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.a != nil {
				if tt.a.Balance != tt.wantBalance {
					t.Fatalf("balance=%d want %d", tt.a.Balance, tt.wantBalance)
				}
				if tt.a.ID != before.ID || tt.a.Owner != before.Owner || tt.a.Active != before.Active || tt.a.Currency != before.Currency {
					t.Fatalf("unrelated fields changed: before=%#v after=%#v", before, *tt.a)
				}
			}
		})
	}
}

func TestAccountWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		a           *Account
		amount      int
		wantErr     error
		wantBalance int
	}{
		{"nil", nil, 10, ErrNilAccount, 0},
		{"inactive", &Account{Balance: 100}, 10, ErrInactive, 100},
		{"zero", &Account{Balance: 100, Active: true}, 0, ErrInvalidAmount, 100},
		{"negative", &Account{Balance: 100, Active: true}, -1, ErrInvalidAmount, 100},
		{"insufficient", &Account{Balance: 100, Active: true}, 101, ErrInsufficientFunds, 100},
		{"exact", &Account{Balance: 100, Active: true}, 100, nil, 0},
		{"partial", &Account{Balance: 100, Active: true}, 40, nil, 60},
		{"zero-balance", &Account{Balance: 0, Active: true}, 1, ErrInsufficientFunds, 0},
		{"negative-balance", &Account{Balance: -1, Active: true}, 1, ErrInsufficientFunds, -1},
		{"large", &Account{Balance: 2_000_000, Active: true}, 1_000_000, nil, 1_000_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.Withdraw(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.a != nil && tt.a.Balance != tt.wantBalance {
				t.Fatalf("balance=%d want %d", tt.a.Balance, tt.wantBalance)
			}
		})
	}
}

func TestAccountTransferTo(t *testing.T) {
	tests := []struct {
		name                 string
		source, destination  *Account
		amount               int
		wantErr              error
		wantSource, wantDest int
	}{
		{"nil-source", nil, &Account{Balance: 10, Active: true}, 1, ErrNilAccount, 0, 10},
		{"nil-destination", &Account{Balance: 10, Active: true}, nil, 1, ErrNilDestination, 10, 0},
		{"same", func() *Account { a := &Account{Balance: 10, Active: true}; return a }(), nil, 1, ErrSameAccount, 10, 10},
		{"source-inactive", &Account{Balance: 10}, &Account{Balance: 1, Active: true}, 1, ErrInactive, 10, 1},
		{"destination-inactive", &Account{Balance: 10, Active: true}, &Account{Balance: 1}, 1, ErrInactive, 10, 1},
		{"zero", &Account{Balance: 10, Active: true}, &Account{Balance: 1, Active: true}, 0, ErrInvalidAmount, 10, 1},
		{"negative", &Account{Balance: 10, Active: true}, &Account{Balance: 1, Active: true}, -1, ErrInvalidAmount, 10, 1},
		{"insufficient", &Account{Balance: 10, Active: true}, &Account{Balance: 1, Active: true}, 11, ErrInsufficientFunds, 10, 1},
		{"exact", &Account{Balance: 10, Active: true}, &Account{Balance: 1, Active: true}, 10, nil, 0, 11},
		{"partial", &Account{Balance: 10, Active: true}, &Account{Balance: 1, Active: true}, 4, nil, 6, 5},
	}
	for i := range tests {
		tt := &tests[i]
		if tt.name == "same" {
			tt.destination = tt.source
		}
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.TransferTo(tt.destination, tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.source != nil && tt.source.Balance != tt.wantSource {
				t.Fatalf("source=%d want %d", tt.source.Balance, tt.wantSource)
			}
			if tt.destination != nil && tt.destination.Balance != tt.wantDest {
				t.Fatalf("destination=%d want %d", tt.destination.Balance, tt.wantDest)
			}
		})
	}
}

func TestAccountActivate(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(nameForTest("case", i), func(t *testing.T) {
			if i == 0 {
				var a *Account
				if !errors.Is(a.Activate(), ErrNilAccount) {
					t.Fatal("nil receiver must return ErrNilAccount")
				}
				return
			}
			a := &Account{ID: i, Active: i%2 == 0, Balance: i * 10}
			if err := a.Activate(); err != nil || !a.Active {
				t.Fatalf("err=%v account=%#v", err, a)
			}
		})
	}
}

func TestAccountDeactivate(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(nameForTest("case", i), func(t *testing.T) {
			if i == 0 {
				var a *Account
				if !errors.Is(a.Deactivate(), ErrNilAccount) {
					t.Fatal("nil receiver must return ErrNilAccount")
				}
				return
			}
			a := &Account{ID: i, Active: i%2 == 0, Balance: i * 10}
			if err := a.Deactivate(); err != nil || a.Active {
				t.Fatalf("err=%v account=%#v", err, a)
			}
		})
	}
}

func TestAccountSnapshot(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(nameForTest("case", i), func(t *testing.T) {
			a := Account{ID: i, Owner: nameForTest("owner", i), Balance: i * 10, Active: i%2 == 0, Currency: "RUB"}
			snapshot := a.Snapshot()
			if !reflect.DeepEqual(snapshot, a) {
				t.Fatalf("snapshot=%#v account=%#v", snapshot, a)
			}
			a.Owner = "changed"
			a.Balance++
			if snapshot.Owner == a.Owner || snapshot.Balance == a.Balance {
				t.Fatal("snapshot changed with source")
			}
		})
	}
}

func TestAccountReset(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(nameForTest("case", i), func(t *testing.T) {
			if i == 0 {
				var a *Account
				a.Reset()
				return
			}
			a := &Account{ID: i, Owner: "A", Balance: i, Active: true, Currency: "RUB"}
			a.Reset()
			if !reflect.DeepEqual(*a, Account{}) {
				t.Fatalf("not reset: %#v", a)
			}
		})
	}
}

func TestAccountApplyBonus(t *testing.T) {
	tests := []struct {
		name        string
		a           *Account
		percent     int
		wantErr     error
		wantBalance int
	}{
		{"nil", nil, 10, ErrNilAccount, 0},
		{"inactive", &Account{Balance: 100}, 10, ErrInactive, 100},
		{"negative", &Account{Balance: 100, Active: true}, -1, ErrInvalidAmount, 100},
		{"zero", &Account{Balance: 100, Active: true}, 0, nil, 100},
		{"ten", &Account{Balance: 100, Active: true}, 10, nil, 110},
		{"hundred", &Account{Balance: 100, Active: true}, 100, nil, 200},
		{"round-down", &Account{Balance: 99, Active: true}, 10, nil, 108},
		{"zero-balance", &Account{Balance: 0, Active: true}, 50, nil, 0},
		{"negative-balance", &Account{Balance: -100, Active: true}, 10, nil, -110},
		{"large", &Account{Balance: 1_000_000, Active: true}, 25, nil, 1_250_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.ApplyBonus(tt.percent)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.a != nil && tt.a.Balance != tt.wantBalance {
				t.Fatalf("balance=%d want %d", tt.a.Balance, tt.wantBalance)
			}
		})
	}
}

func TestAccountSameOwner(t *testing.T) {
	tests := []struct {
		name        string
		left, right string
		want        bool
	}{
		{"empty", "", "", true},
		{"same", "A", "A", true},
		{"different", "A", "B", false},
		{"case", "Go", "go", false},
		{"space", "A", "A ", false},
		{"both-spaces", "  ", "  ", true},
		{"unicode-same", "Мария", "Мария", true},
		{"unicode-different", "Мария", "мария", false},
		{"long", "very long owner", "very long owner", true},
		{"punctuation", "A-B", "A_B", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Account{Owner: tt.left}).SameOwner(Account{Owner: tt.right}); got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
		})
	}
}

func TestCartAdd(t *testing.T) {
	tests := []struct {
		name    string
		cart    *Cart
		sku     string
		count   int
		wantErr error
		want    map[string]int
	}{
		{"nil", nil, "a", 1, ErrNilCart, nil},
		{"empty-sku", &Cart{MaxItems: 5}, "", 1, ErrInvalidSKU, nil},
		{"zero", &Cart{MaxItems: 5}, "a", 0, ErrInvalidAmount, nil},
		{"negative", &Cart{MaxItems: 5}, "a", -1, ErrInvalidAmount, nil},
		{"zero-limit", &Cart{MaxItems: 0}, "a", 1, ErrCartLimit, nil},
		{"over-limit", &Cart{MaxItems: 2}, "a", 3, ErrCartLimit, nil},
		{"init-map", &Cart{MaxItems: 2}, "a", 2, nil, map[string]int{"a": 2}},
		{"append-existing", &Cart{Items: map[string]int{"a": 1}, MaxItems: 3}, "a", 2, nil, map[string]int{"a": 3}},
		{"different-sku", &Cart{Items: map[string]int{"a": 1}, MaxItems: 3}, "b", 2, nil, map[string]int{"a": 1, "b": 2}},
		{"existing-total-limit", &Cart{Items: map[string]int{"a": 2}, MaxItems: 3}, "b", 2, ErrCartLimit, map[string]int{"a": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cart.Add(tt.sku, tt.count)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v want %v", err, tt.wantErr)
			}
			if tt.cart != nil && !reflect.DeepEqual(tt.cart.Items, tt.want) {
				t.Fatalf("items=%#v want %#v", tt.cart.Items, tt.want)
			}
		})
	}
}

func TestCartTotalItems(t *testing.T) {
	tests := []struct {
		name  string
		items map[string]int
		want  int
	}{
		{"nil", nil, 0},
		{"empty", map[string]int{}, 0},
		{"one", map[string]int{"a": 1}, 1},
		{"two", map[string]int{"a": 1, "b": 2}, 3},
		{"zero", map[string]int{"a": 0}, 0},
		{"negative", map[string]int{"a": -1}, -1},
		{"mixed", map[string]int{"a": 3, "b": -1}, 2},
		{"empty-sku", map[string]int{"": 5}, 5},
		{"large", map[string]int{"a": 1_000_000}, 1_000_000},
		{"many", map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Cart{Items: tt.items}).TotalItems(); got != tt.want {
				t.Fatalf("got %d want %d", got, tt.want)
			}
		})
	}
}

func nameForTest(prefix string, i int) string {
	return prefix + "-" + string(rune('a'+i))
}
