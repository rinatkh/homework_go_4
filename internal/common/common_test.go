package common

import (
	"reflect"
	"testing"

	"github.com/rinatkh/homework_go_4/internal/methods"
)

func TestApplyOperations(t *testing.T) {
	tests := []struct {
		name     string
		accounts map[int]*methods.Account
		ops      []Operation
		balances map[int]int
		errors   []string
	}{
		{"empty", map[int]*methods.Account{}, nil, map[int]int{}, []string{}},
		{"deposit", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}}, []Operation{{1, "deposit", 20}}, map[int]int{1: 120}, []string{}},
		{"withdraw", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}}, []Operation{{1, "withdraw", 20}}, map[int]int{1: 80}, []string{}},
		{"inactive", map[int]*methods.Account{1: {ID: 1, Balance: 100}}, []Operation{{1, "deposit", 20}}, map[int]int{1: 100}, []string{"account 1: inactive account"}},
		{"missing", map[int]*methods.Account{}, []Operation{{1, "deposit", 20}}, map[int]int{}, []string{"account 1: not found"}},
		{"nil-account", map[int]*methods.Account{1: nil}, []Operation{{1, "deposit", 20}}, map[int]int{}, []string{"account 1: not found"}},
		{"unknown", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}}, []Operation{{1, "refund", 20}}, map[int]int{1: 100}, []string{"account 1: unknown operation \"refund\""}},
		{"invalid-amount", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}}, []Operation{{1, "deposit", 0}}, map[int]int{1: 100}, []string{"account 1: invalid amount"}},
		{"insufficient", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}}, []Operation{{1, "withdraw", 101}}, map[int]int{1: 100}, []string{"account 1: insufficient funds"}},
		{"continue", map[int]*methods.Account{1: {ID: 1, Balance: 100, Active: true}, 2: {ID: 2, Balance: 10, Active: true}}, []Operation{{1, "deposit", 20}, {3, "deposit", 1}, {2, "withdraw", 5}}, map[int]int{1: 120, 2: 5}, []string{"account 3: not found"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyOperations(tt.accounts, tt.ops)
			if got.Balances == nil || got.Errors == nil {
				t.Fatalf("result fields must be initialized: %#v", got)
			}
			if !reflect.DeepEqual(got.Balances, tt.balances) {
				t.Fatalf("balances=%#v want %#v", got.Balances, tt.balances)
			}
			if !reflect.DeepEqual(got.Errors, tt.errors) {
				t.Fatalf("errors=%#v want %#v", got.Errors, tt.errors)
			}
		})
	}
}

func TestBuildAccountIndex(t *testing.T) {
	cases := [][]methods.Account{
		nil,
		{},
		{{ID: 1, Owner: "A"}},
		{{ID: 1, Owner: "A"}, {ID: 2, Owner: "B"}},
		{{ID: 1, Owner: "A"}, {ID: 1, Owner: "B"}},
		{{ID: 0, Owner: "Zero"}},
		{{ID: -1, Owner: "N"}},
		{{ID: 1, Balance: 100, Active: true}},
		{{ID: 1, Currency: "RUB"}},
		{{ID: 1, Owner: "A"}, {ID: 2, Owner: "B"}, {ID: 3, Owner: "C"}},
	}
	for i, source := range cases {
		t.Run(commonTestName(i), func(t *testing.T) {
			before := cloneAccountsForTest(source)
			got := BuildAccountIndex(source)
			if got == nil {
				t.Fatal("result map must be initialized")
			}
			for id, account := range got {
				if account == nil || account.ID != id {
					t.Fatalf("bad account for id %d: %#v", id, account)
				}
			}
			if len(got) > 1 {
				seen := map[*methods.Account]bool{}
				for _, account := range got {
					if seen[account] {
						t.Fatal("different ids point to the same account")
					}
					seen[account] = true
				}
			}
			for _, account := range got {
				account.Owner = "changed"
				break
			}
			if !reflect.DeepEqual(source, before) {
				t.Fatalf("source changed: %#v", source)
			}
		})
	}
}

func TestApplyUserEvents(t *testing.T) {
	yes := true
	no := false
	tests := []struct {
		name       string
		users      map[int]string
		active     map[int]bool
		events     []Event
		wantUsers  map[int]string
		wantActive map[int]bool
	}{
		{"nil", nil, nil, nil, map[int]string{}, map[int]bool{}},
		{"empty", map[int]string{}, map[int]bool{}, []Event{}, map[int]string{}, map[int]bool{}},
		{"name-new", nil, nil, []Event{{UserID: 1, Name: "A"}}, map[int]string{1: "A"}, map[int]bool{}},
		{"active-new", nil, nil, []Event{{UserID: 1, Active: &yes}}, map[int]string{}, map[int]bool{1: true}},
		{"both-new", nil, nil, []Event{{UserID: 1, Name: "A", Active: &no}}, map[int]string{1: "A"}, map[int]bool{1: false}},
		{"ignore-empty", map[int]string{1: "A"}, map[int]bool{1: true}, []Event{{UserID: 1}}, map[int]string{1: "A"}, map[int]bool{1: true}},
		{"update-name", map[int]string{1: "A"}, map[int]bool{1: false}, []Event{{UserID: 1, Name: "B"}}, map[int]string{1: "B"}, map[int]bool{1: false}},
		{"update-active", map[int]string{1: "A"}, map[int]bool{1: false}, []Event{{UserID: 1, Active: &yes}}, map[int]string{1: "A"}, map[int]bool{1: true}},
		{"order", nil, nil, []Event{{UserID: 1, Name: "A", Active: &yes}, {UserID: 1, Name: "B", Active: &no}}, map[int]string{1: "B"}, map[int]bool{1: false}},
		{"many", map[int]string{1: "A"}, map[int]bool{1: false}, []Event{{UserID: 2, Name: "B", Active: &yes}, {UserID: 1, Name: "AA"}}, map[int]string{1: "AA", 2: "B"}, map[int]bool{1: false, 2: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usersBefore := tt.users
			activeBefore := tt.active
			gotUsers, gotActive := ApplyUserEvents(tt.users, tt.active, tt.events)
			if gotUsers == nil || gotActive == nil {
				t.Fatalf("returned maps must be initialized: %#v %#v", gotUsers, gotActive)
			}
			if !reflect.DeepEqual(gotUsers, tt.wantUsers) || !reflect.DeepEqual(gotActive, tt.wantActive) {
				t.Fatalf("got users=%#v active=%#v; want users=%#v active=%#v", gotUsers, gotActive, tt.wantUsers, tt.wantActive)
			}
			if usersBefore != nil && reflect.ValueOf(usersBefore).Pointer() != reflect.ValueOf(gotUsers).Pointer() {
				t.Fatal("non-nil users map must be updated in place")
			}
			if activeBefore != nil && reflect.ValueOf(activeBefore).Pointer() != reflect.ValueOf(gotActive).Pointer() {
				t.Fatal("non-nil active map must be updated in place")
			}
		})
	}
}

func cloneAccountsForTest(source []methods.Account) []methods.Account {
	if source == nil {
		return nil
	}
	result := make([]methods.Account, len(source))
	copy(result, source)
	return result
}

func commonTestName(i int) string {
	return "case-" + string(rune('a'+i))
}
