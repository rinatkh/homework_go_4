package integration

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/rinatkh/homework_go_4/internal/structs"
)

func TestCommands(t *testing.T) {
	bad, good, saved := structs.LayoutSizes()
	cases := map[string]string{
		"01_maps": "go=2 missing=42\n" +
			"moscow=2 tags=go:2 api:1 sql:1\n" +
			"low=[book] value=300",
		"02_methods": "account=#1 Maria: 1200 RUB\n" +
			"original=Maria copy=Masha\n" +
			"cart=3",
		"03_structs": fmt.Sprintf(
			"user=1:Masha logins=1\n"+
				"admin=admin:2:Ada read=true\n"+
				"value=%d pointer=%d layout=%d/%d saved=%d",
			structs.UserValueSize(structs.User{}), structs.UserPointerSize(&structs.User{}), bad, good, saved,
		),
		"04_common": "balances=map[1:1200 2:100] errors=1\n" +
			"users=map[1:Maria] active=map[1:true]",
	}

	for command, want := range cases {
		command, want := command, want
		t.Run(command, func(t *testing.T) {
			out, err := exec.Command("go", "run", "../../cmd/"+command).CombinedOutput()
			if err != nil {
				t.Fatalf("go run ./cmd/%s failed: %v\n%s", command, err, out)
			}

			got := strings.TrimSpace(string(out))
			if got != want {
				t.Fatalf("unexpected output for %s:\n\ngot:\n%s\n\nwant:\n%s", command, got, want)
			}
		})
	}
}
