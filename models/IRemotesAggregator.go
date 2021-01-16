package models

const (
	ErrorNoSuchRef = Error("No such reference")
)

type IRemotesAggregator interface {
	GetRemoteRef(remoteName string) (string, error)
}
