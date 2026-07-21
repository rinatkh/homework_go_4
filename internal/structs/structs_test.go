package structs

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestUpdateUserName(t *testing.T) {
	tests := []struct {
		name  string
		users map[int]User
		id    int
		value string
		want  bool
	}{
		{"nil", nil, 1, "A", false},
		{"empty", map[int]User{}, 1, "A", false},
		{"missing", map[int]User{1: {ID: 1, Name: "A"}}, 2, "B", false},
		{"basic", map[int]User{1: {ID: 1, Name: "A"}}, 1, "B", true},
		{"empty-name", map[int]User{1: {ID: 1, Name: "A"}}, 1, "", true},
		{"same-name", map[int]User{1: {ID: 1, Name: "A"}}, 1, "A", true},
		{"zero-id", map[int]User{0: {ID: 0, Name: "A"}}, 0, "B", true},
		{"negative-id", map[int]User{-1: {ID: -1, Name: "A"}}, -1, "B", true},
		{"unicode", map[int]User{1: {ID: 1, Name: "A"}}, 1, "Мария", true},
		{"keep-fields", map[int]User{1: {ID: 1, Name: "A", Active: true, LoginCount: 3, Tags: []string{"go"}}}, 1, "B", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeLen := len(tt.users)
			got := UpdateUserName(tt.users, tt.id, tt.value)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if got && tt.users[tt.id].Name != tt.value {
				t.Fatalf("name=%q want %q", tt.users[tt.id].Name, tt.value)
			}
			if !got && len(tt.users) != beforeLen {
				t.Fatal("missing user was added")
			}
		})
	}
}

func TestActivateUser(t *testing.T) {
	tests := []struct {
		name  string
		users map[int]User
		id    int
		want  bool
	}{
		{"nil", nil, 1, false},
		{"empty", map[int]User{}, 1, false},
		{"missing", map[int]User{1: {ID: 1}}, 2, false},
		{"inactive", map[int]User{1: {ID: 1}}, 1, true},
		{"already-active", map[int]User{1: {ID: 1, Active: true}}, 1, true},
		{"zero-id", map[int]User{0: {ID: 0}}, 0, true},
		{"negative-id", map[int]User{-1: {ID: -1}}, -1, true},
		{"keep-name", map[int]User{1: {ID: 1, Name: "A"}}, 1, true},
		{"keep-logins", map[int]User{1: {ID: 1, LoginCount: 7}}, 1, true},
		{"keep-tags", map[int]User{1: {ID: 1, Tags: []string{"go"}}}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActivateUser(tt.users, tt.id)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if got && !tt.users[tt.id].Active {
				t.Fatal("user is not active")
			}
		})
	}
}

func TestIncrementLoginCount(t *testing.T) {
	tests := []struct {
		name  string
		users map[int]User
		id    int
		want  bool
		count int
	}{
		{"nil", nil, 1, false, 0},
		{"empty", map[int]User{}, 1, false, 0},
		{"missing", map[int]User{1: {LoginCount: 3}}, 2, false, 0},
		{"zero", map[int]User{1: {LoginCount: 0}}, 1, true, 1},
		{"positive", map[int]User{1: {LoginCount: 5}}, 1, true, 6},
		{"negative", map[int]User{1: {LoginCount: -2}}, 1, true, -1},
		{"zero-id", map[int]User{0: {LoginCount: 1}}, 0, true, 2},
		{"negative-id", map[int]User{-1: {LoginCount: 1}}, -1, true, 2},
		{"keep-active", map[int]User{1: {Active: true, LoginCount: 1}}, 1, true, 2},
		{"large", map[int]User{1: {LoginCount: 1_000_000}}, 1, true, 1_000_001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IncrementLoginCount(tt.users, tt.id)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if got && tt.users[tt.id].LoginCount != tt.count {
				t.Fatalf("count=%d want %d", tt.users[tt.id].LoginCount, tt.count)
			}
		})
	}
}

func TestCopyUsers(t *testing.T) {
	cases := []map[int]User{
		nil,
		{},
		{1: {ID: 1}},
		{1: {ID: 1, Name: "A"}, 2: {ID: 2, Name: "B"}},
		{0: {ID: 0, Tags: []string{}}},
		{-1: {ID: -1, Tags: []string{"go"}}},
		{1: {ID: 1, Active: true}},
		{1: {ID: 1, LoginCount: 10}},
		{1: {ID: 1, Tags: []string{"go", "api"}}},
		{1: {ID: 1, Tags: nil}, 2: {ID: 2, Tags: []string{"sql"}}},
	}
	for i, source := range cases {
		t.Run(testName(i), func(t *testing.T) {
			got := CopyUsers(source)
			if got == nil {
				t.Fatal("result map must be initialized")
			}
			if len(got) != len(source) {
				t.Fatalf("len(got)=%d want %d", len(got), len(source))
			}
			for id, wantUser := range source {
				if gotUser, ok := got[id]; !ok || !reflect.DeepEqual(gotUser, wantUser) {
					t.Fatalf("user %d: got %#v want %#v", id, gotUser, wantUser)
				}
			}
			if len(got) > 0 {
				for id, user := range got {
					user.Name = "changed"
					if user.Tags != nil {
						user.Tags = append(user.Tags, "new")
					}
					got[id] = user
					if source[id].Name == "changed" {
						t.Fatal("map values are not independent")
					}
					if len(source[id].Tags) > 0 && len(got[id].Tags) == len(source[id].Tags) {
						t.Fatal("tags slice is not independent")
					}
					break
				}
			}
		})
	}
}

func TestBuildPointerIndex(t *testing.T) {
	cases := [][]User{
		nil,
		{},
		{{ID: 1, Name: "A"}},
		{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		{{ID: 1, Name: "A"}, {ID: 1, Name: "B"}},
		{{ID: 0, Name: "Zero"}},
		{{ID: -1, Name: "N"}},
		{{ID: 1, Tags: []string{"go"}}},
		{{ID: 1, Active: true, LoginCount: 3}},
		{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}, {ID: 3, Name: "C"}},
	}
	for i, source := range cases {
		t.Run(testName(i), func(t *testing.T) {
			before := cloneUsersForTest(source)
			got := BuildPointerIndex(source)
			if got == nil {
				t.Fatal("result map must be initialized")
			}
			for id, ptr := range got {
				if ptr == nil || ptr.ID != id {
					t.Fatalf("bad pointer for id %d: %#v", id, ptr)
				}
			}
			if len(got) > 1 {
				seen := map[*User]bool{}
				for _, ptr := range got {
					if seen[ptr] {
						t.Fatal("different ids point to the same object")
					}
					seen[ptr] = true
				}
			}
			for _, ptr := range got {
				ptr.Name = "changed"
				break
			}
			if !reflect.DeepEqual(source, before) {
				t.Fatalf("source slice changed: %#v", source)
			}
		})
	}
}

func TestRenameThroughPointer(t *testing.T) {
	user := func(name string) *User { return &User{Name: name} }
	tests := []struct {
		name  string
		users map[int]*User
		id    int
		value string
		want  bool
	}{
		{"nil-map", nil, 1, "B", false},
		{"empty", map[int]*User{}, 1, "B", false},
		{"missing", map[int]*User{1: user("A")}, 2, "B", false},
		{"nil-value", map[int]*User{1: nil}, 1, "B", false},
		{"basic", map[int]*User{1: user("A")}, 1, "B", true},
		{"empty-name", map[int]*User{1: user("A")}, 1, "", true},
		{"same-name", map[int]*User{1: user("A")}, 1, "A", true},
		{"zero-id", map[int]*User{0: user("A")}, 0, "B", true},
		{"negative-id", map[int]*User{-1: user("A")}, -1, "B", true},
		{"unicode", map[int]*User{1: user("A")}, 1, "Мария", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenameThroughPointer(tt.users, tt.id, tt.value)
			if got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
			if got && tt.users[tt.id].Name != tt.value {
				t.Fatalf("name=%q want %q", tt.users[tt.id].Name, tt.value)
			}
		})
	}
}

func TestUserDisplayName(t *testing.T) {
	tests := []struct {
		user User
		want string
	}{
		{User{}, "0:"},
		{User{ID: 1, Name: "A"}, "1:A"},
		{User{ID: -1, Name: "A"}, "-1:A"},
		{User{ID: 0, Name: "Zero"}, "0:Zero"},
		{User{ID: 10, Name: "Maria"}, "10:Maria"},
		{User{ID: 2, Name: "A B"}, "2:A B"},
		{User{ID: 3, Name: "Мария"}, "3:Мария"},
		{User{ID: 4, Name: ""}, "4:"},
		{User{ID: 5, Name: "  "}, "5:  "},
		{User{ID: 999, Name: "User"}, "999:User"},
	}
	for i, tt := range tests {
		t.Run(testName(i), func(t *testing.T) {
			if got := tt.user.DisplayName(); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestAdminDisplayName(t *testing.T) {
	tests := []struct {
		admin Admin
		want  string
	}{
		{Admin{}, "admin:0:"},
		{Admin{User: User{ID: 1, Name: "A"}}, "admin:1:A"},
		{Admin{User: User{ID: -1, Name: "A"}}, "admin:-1:A"},
		{Admin{User: User{ID: 0, Name: "Zero"}}, "admin:0:Zero"},
		{Admin{User: User{ID: 10, Name: "Maria"}}, "admin:10:Maria"},
		{Admin{User: User{ID: 2, Name: "A B"}}, "admin:2:A B"},
		{Admin{User: User{ID: 3, Name: "Мария"}}, "admin:3:Мария"},
		{Admin{User: User{ID: 4}}, "admin:4:"},
		{Admin{User: User{ID: 5, Name: "  "}}, "admin:5:  "},
		{Admin{User: User{ID: 999, Name: "Admin"}, Permissions: []string{"read"}}, "admin:999:Admin"},
	}
	for i, tt := range tests {
		t.Run(testName(i), func(t *testing.T) {
			if got := tt.admin.DisplayName(); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestAdminHasPermission(t *testing.T) {
	tests := []struct {
		permissions []string
		permission  string
		want        bool
	}{
		{nil, "read", false},
		{[]string{}, "read", false},
		{[]string{"read"}, "read", true},
		{[]string{"read"}, "write", false},
		{[]string{"read", "write"}, "write", true},
		{[]string{"read", "read"}, "read", true},
		{[]string{"Read"}, "read", false},
		{[]string{" read"}, "read", false},
		{[]string{""}, "", true},
		{[]string{"бан"}, "бан", true},
	}
	for i, tt := range tests {
		t.Run(testName(i), func(t *testing.T) {
			if got := (Admin{Permissions: tt.permissions}).HasPermission(tt.permission); got != tt.want {
				t.Fatalf("got %t want %t", got, tt.want)
			}
		})
	}
}

func TestRenameEmbeddedUser(t *testing.T) {
	names := []string{"A", "", "B", "Maria", "Мария", "A B", "  ", "Go", "admin-1", "long name"}
	for i, name := range names {
		t.Run(testName(i), func(t *testing.T) {
			if i == 0 {
				RenameEmbeddedUser(nil, name)
				return
			}
			admin := &Admin{User: User{ID: i, Name: "old"}, Permissions: []string{"read"}}
			RenameEmbeddedUser(admin, name)
			if admin.Name != name || admin.ID != i || !reflect.DeepEqual(admin.Permissions, []string{"read"}) {
				t.Fatalf("admin=%#v", admin)
			}
		})
	}
}

func TestPremiumAccountTotal(t *testing.T) {
	tests := []struct {
		balance int
		bonus   int
		want    int
	}{
		{0, 0, 0},
		{100, 0, 100},
		{0, 10, 10},
		{100, 10, 110},
		{-100, 10, -90},
		{100, -10, 90},
		{-100, -10, -110},
		{1, 1, 2},
		{1_000_000, 2_000_000, 3_000_000},
		{99, 1, 100},
	}
	for i, tt := range tests {
		t.Run(testName(i), func(t *testing.T) {
			p := PremiumAccount{Account: Account{Balance: tt.balance}, Bonus: tt.bonus}
			if got := p.Total(); got != tt.want {
				t.Fatalf("got %d want %d", got, tt.want)
			}
		})
	}
}

func TestPremiumLabel(t *testing.T) {
	tests := []struct {
		p    PremiumAccount
		want string
	}{
		{PremiumAccount{}, "0:0 total=0"},
		{PremiumAccount{Account: Account{ID: 1}}, "1:0 total=0"},
		{PremiumAccount{Account: Account{ID: 1, Balance: 100}}, "1:100 total=100"},
		{PremiumAccount{Account: Account{ID: 1, Balance: 100}, Bonus: 10}, "1:100 total=110"},
		{PremiumAccount{Account: Account{ID: -1, Balance: 100}, Bonus: 10}, "-1:100 total=110"},
		{PremiumAccount{Account: Account{ID: 2, Balance: -100}, Bonus: 10}, "2:-100 total=-90"},
		{PremiumAccount{Account: Account{ID: 3, Balance: 0}, Bonus: -10}, "3:0 total=-10"},
		{PremiumAccount{Account: Account{ID: 4, Balance: 1}, Bonus: 1}, "4:1 total=2"},
		{PremiumAccount{Account: Account{ID: 5, Balance: 1_000_000}, Bonus: 2_000_000}, "5:1000000 total=3000000"},
		{PremiumAccount{Account: Account{ID: 999, Balance: 99}, Bonus: 1}, "999:99 total=100"},
	}
	for i, tt := range tests {
		t.Run(testName(i), func(t *testing.T) {
			if got := PremiumLabel(tt.p); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestUserValueSize(t *testing.T) {
	users := []User{
		{}, {ID: 1}, {Name: "A"}, {Active: true}, {LoginCount: 10},
		{Tags: nil}, {Tags: []string{}}, {Tags: []string{"go"}},
		{ID: -1, Name: "Мария", Active: true, LoginCount: 7},
		{ID: 999, Name: "long name", Tags: []string{"go", "api"}},
	}
	want := unsafe.Sizeof(User{})
	for i, user := range users {
		t.Run(testName(i), func(t *testing.T) {
			if got := UserValueSize(user); got != want {
				t.Fatalf("got %d want %d", got, want)
			}
		})
	}
}

func TestUserPointerSize(t *testing.T) {
	users := []*User{
		nil, {}, {ID: 1}, {Name: "A"}, {Active: true},
		{LoginCount: 10}, {Tags: nil}, {Tags: []string{}},
		{ID: -1, Name: "Мария"}, {ID: 999, Tags: []string{"go"}},
	}
	var pointer *User
	want := unsafe.Sizeof(pointer)
	for i, user := range users {
		t.Run(testName(i), func(t *testing.T) {
			if got := UserPointerSize(user); got != want {
				t.Fatalf("got %d want %d", got, want)
			}
		})
	}
}

func TestLayoutSizes(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(testName(i), func(t *testing.T) {
			bad, good, saved := LayoutSizes()
			if bad != unsafe.Sizeof(BadLayout{}) || good != unsafe.Sizeof(GoodLayout{}) {
				t.Fatalf("bad=%d good=%d", bad, good)
			}
			if bad <= good || saved != bad-good {
				t.Fatalf("bad=%d good=%d saved=%d", bad, good, saved)
			}
		})
	}
}

func cloneUsersForTest(source []User) []User {
	if source == nil {
		return nil
	}
	result := make([]User, len(source))
	copy(result, source)
	return result
}

func testName(i int) string {
	return "case-" + string(rune('a'+i))
}
