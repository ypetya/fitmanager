package models

type IRemote interface {
	GetType() RemoteType
	GetRef() string
}
