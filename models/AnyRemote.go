package models

import "time"

// An union remote type for storing data on disk:
//
// min cut between serialization and different types
//
// can point to an exercise, directory
// might have an external Id
// or mark a remote in general
// TODO split these up and make it more specific
type AnyRemote struct {
	Name     string     `json:"name,omitempty"`
	Target   RemoteType `json:"target,omitempty"`
	File     *File      `json:"file,omitempty"`
	Id       string     `json:"id,omitempty"`
	LastSync int64      `json:"u,omitempty"`
}

// Instantiates a new Directory type remote, filled out the File
func NewRemoteDirectory(name string, dir string) AnyRemote {
	file := newFromDir(dir)
	return AnyRemote{
		Name:   name,
		Target: Directory,
		File:   &file,
	}
}

func (r AnyRemote) GetType() RemoteType {
	return r.Target
}

// Ref returns either the File.Path as id or the stored remoteId
// filepath is filled-out for directory type remotes and localDB
// and contains FileName.
// (The path of the filename is provided by the data-store's remotes context)
func (r AnyRemote) GetRef() string {
	t := r.GetType()
	switch {
	case t == Directory || t == LocalDB:
		if r.File == nil {
			panic("Remote missing mandatory File attribute:" + t)
		}
		return r.File.Path
	case t == GarminConnect:
		return r.Id
	default:
		return ""
	}
}

func (r AnyRemote) Fill(ref string) AnyRemote {
	t := r.GetType()
	switch {
	case t == Directory || t == LocalDB:
		r.File = &File{Path: ref, LastSeen: time.Now().Unix()}
	case t == GarminConnect:
		r.Id = ref
		r.LastSync = time.Now().Unix()
	}

	return r
}

// This method use to check equality in the same exercise context:
// * remotes are equal if both of them are localDB references
// * or for the same remote(remote type equals and name equals)
func (r AnyRemote) equals(o IEquals) bool {
	if v, ok := o.(AnyRemote); ok {
		if r.Target == v.Target && r.Target == LocalDB {
			return true
		}
		return r.Name == v.Name && r.Target == v.Target
	}
	return false
}
