package models

type RemoteType string

const (
	Directory     RemoteType = "Folder"
	LocalDB       RemoteType = "local"
	GarminConnect RemoteType = "GarminConnect"
)
