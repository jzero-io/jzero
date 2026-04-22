package hooks

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunPrintsHookCommandBeforeExecution(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("quiet", false, "")

	oldRunOutputFn := runOutputFn
	oldPrintHookCommandFn := printHookCommandFn
	oldPrintHookOutputFn := printHookOutputFn
	t.Cleanup(func() {
		runOutputFn = oldRunOutputFn
		printHookCommandFn = oldPrintHookCommandFn
		printHookOutputFn = oldPrintHookOutputFn
	})

	var calls []string
	printHookCommandFn = func(command string) {
		calls = append(calls, "print:"+command)
	}
	runOutputFn = func(arg, dir string, env ...string) (string, error) {
		calls = append(calls, "run:"+arg)
		return "hook output", nil
	}
	printHookOutputFn = func(output string) {
		calls = append(calls, "output:"+output)
	}

	if err := Run(cmd, "Before", "gen", []string{"echo test"}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := []string{
		"print:echo test",
		"run:echo test",
		"output:hook output",
	}

	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("Run() calls = %v, want %v", calls, want)
	}
}
