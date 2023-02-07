package differ

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/leep-frog/command"
	"github.com/leep-frog/command/sourcerer"
)

type Differ struct{}

func CLI() sourcerer.CLI {
	return &Differ{}
}

func (d *Differ) Name() string    { return "df" }
func (d *Differ) Changed() bool   { return false }
func (d *Differ) Setup() []string { return nil }

const (
	splitString = "> but was: <"
)

func (d *Differ) Node() command.Node {
	inputArg := command.Arg("INPUT", "Input to diff", command.Default("bleh"))
	jsonFlag := command.BoolFlag("json", 'j', "Compares json objects if set to true")
	return command.SerialNodes(
		command.FlagProcessor(
			jsonFlag,
		),
		inputArg,
		&command.ExecutorProcessor{func(o command.Output, d *command.Data) error {
			input := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(inputArg.Get(d)), "<"), ">")
			parts := strings.Split(input, splitString)
			if len(parts) != 2 {
				return o.Stderrf("Input did not contain split string: %q", splitString)
			}

			diff := cmp.Diff(parts[0], parts[1])
			if jsonFlag.Get(d) {
				var a, b interface{}
				if err := json.Unmarshal([]byte(a.(string)), &a); err != nil {
					o.Stderrf("Failed to parse LHS json: %v", err)
				}
				if err := json.Unmarshal([]byte(b.(string)), &b); err != nil {
					o.Stderrf("Failed to parse RHS json: %v", err)
				}
				diff = cmp.Diff(a, b)
			}

			if diff != "" {
				// Only send diff to stdout
				err := o.Stderrln("Objects do not match (-lhs, +rhs):")
				o.Stdoutln(diff)
				fmt.Println(diff)
				return err
			}

			// Only diff, or nothing is sent to stdout
			o.Stderrln("Objects match!")
			return nil
		}},
	)
}
