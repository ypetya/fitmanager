package console

import "testing"

var (
	calls  = []string{}
	params = [][]string{}
	ce     = Command{
		executors: map[string]CommandExecutor{
			"command": func(inputParams []string, _ interface{}) {
				calls = append(calls, "command")
				params = append(params, inputParams)
			},
			"failing": func(inputParams []string, _ interface{}) {
				calls = append(calls, "failing")
				params = append(params, inputParams)
				panic("failing method called")
			},
		},
	}
)

func reset() {
	calls = []string{}
	params = [][]string{}
}

func TestExecuteCommand(t *testing.T) {
	// GIVEN
	reset()
	// WHEN
	ce.Execute("command", []string{})
	// THEN
	if calls[0] != "command" {
		t.Error("Execute should take first arg as a command")
	}
}
func TestExecutePassingArgs(t *testing.T) {
	// GIVEN
	p := []string{"first", "second"}
	reset()
	// WHEN
	ce.Execute("command", p)
	// THEN
	if len(params[0]) != len(p) {
		t.Error("Execute should pass remaining args to executor")
	}
}

func TestExecutingNonExistingCommandReturnsFalse(t *testing.T) {
	// GIVEN
	reset()
	// WHEN
	ret := ce.Execute("funky", []string{})
	// THEN
	if ret {
		t.Error("Execute method should return false when no such command")
	}
}
func TestExecutingExistingCommandReturnsTrue(t *testing.T) {
	// GIVEN
	reset()
	// WHEN
	ret := ce.Execute("funky", []string{})
	// THEN
	if ret {
		t.Error("Execute method should return false when no such command")
	}
}
func TestExecutingFailingCommandReturnsFalse(t *testing.T) {
	// GIVEN
	reset()
	// WHEN
	ret := ce.Execute("failing", []string{})
	// THEN
	if ret {
		t.Error("Execute method should return false when command throws an exception")
	}
}
