package console

type CommandExecutor func([]string, interface{})

type ICommands interface {
	AddCommand(name string, callback CommandExecutor)
	Execute(name string, args []string) bool
	GetSubject() interface{}
}
