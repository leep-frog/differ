package differ

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leep-frog/command"
)

func TestExecute(t *testing.T) {
	for _, test := range []struct {
		name string
		etc  *command.ExecuteTestCase
	}{
		{
			name: "Compares two equivalent strings",
			etc: &command.ExecuteTestCase{
				Args:       []string{` <abc> but was: <abc>  `},
				WantStderr: "Objects match!\n",
			},
		},
		{
			name: "Compares two different strings",
			etc: &command.ExecuteTestCase{
				Args:       []string{` <abc> but was: <def>  `},
				WantStderr: "Objects do not match (-lhs, +rhs):\n",
				WantErr:    fmt.Errorf("Objects do not match (-lhs, +rhs):"),
				WantStdout: cmp.Diff("abc", "def") + "\n",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.etc.Node = CLI().Node()
			test.etc.SkipDataCheck = true
			command.ExecuteTest(t, test.etc)
		})
	}
}
