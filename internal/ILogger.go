package internal

// interface to pass / mock logging of DataStore
type ILogger interface {
	printLn(str string)
}
