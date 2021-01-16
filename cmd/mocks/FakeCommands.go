package mocks

import "github.com/ypetya/fitmanager/console"

type FakeCommands struct {
	Calls    []string
	Names    []string
	Commands map[string]console.CommandExecutor
}

func (f *FakeCommands) AddCommand(name string, callback console.CommandExecutor) {
	f.Calls = append(f.Calls, "AddCommand")
	f.Commands[name] = callback
}

func (f *FakeCommands) Execute(name string, args []string) bool {
	f.Calls = append(f.Calls, "Execute")
	f.Names = append(f.Names, name)
	return true
}

func (f *FakeCommands) GetSubject() interface{} {
	return nil
}

func (f *FakeCommands) HasCommand(command string) bool {
	_, ok := f.Commands[command]
	return ok
}

func NewFakeCommands() FakeCommands {
	return FakeCommands{
		Commands: map[string]console.CommandExecutor{},
	}
}
