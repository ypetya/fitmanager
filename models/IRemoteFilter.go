package models

// RemoteFilter is a special filter
// which has a builder method used for constructing filter criteria
type IRemoteFilter interface {
	Remote(remote string, condition []string)
	Filter(exercises []Exercise) []Exercise
}
