package structs

import (
	"fmt"
	"strings"
)

func Example() string {
	var out strings.Builder

	users := map[int]User{1: {ID: 1, Name: "Maria", Tags: []string{"go"}}}
	_ = UpdateUserName(users, 1, "Masha")
	_ = IncrementLoginCount(users, 1)
	fmt.Fprintf(&out, "user=%s logins=%d\n", users[1].DisplayName(), users[1].LoginCount)

	admin := Admin{User: User{ID: 2, Name: "Ada"}, Permissions: []string{"read", "ban"}}
	fmt.Fprintf(&out, "admin=%s read=%t\n", admin.DisplayName(), admin.HasPermission("read"))

	bad, good, saved := LayoutSizes()
	fmt.Fprintf(&out, "value=%d pointer=%d layout=%d/%d saved=%d",
		UserValueSize(User{}), UserPointerSize(&User{}), bad, good, saved)

	return out.String()
}
