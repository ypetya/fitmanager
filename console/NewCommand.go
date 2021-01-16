package console

func NewCommand(subject interface{}) Command {
	return Command{
		executors: map[string]CommandExecutor{},
		subject:   subject,
	}
}
