package console

import "fmt"

type Command struct {
	executors map[string]CommandExecutor
	subject   interface{}
}

func (c *Command) AddCommand(name string, callback CommandExecutor) {
	c.executors[name] = callback
}

func (c *Command) Execute(name string, args []string) (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed:", r)
			ret = false
		}
	}()
	if fn, ok := c.executors[name]; ok {
		fmt.Println("Executing command:", name, args)
		fn(args, c.subject)
		return true
	} else {
		fmt.Println("Command not found:", name)
		return false
	}
}

func (c *Command) GetSubject() interface{} {
	return c.subject
}
