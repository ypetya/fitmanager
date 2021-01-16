package main

import (
	"fmt"
	"time"

	m "../models"
)

func commandAddRemote(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}
	dsDir := params[0]

	if len(params) < 2 {
		panic("remote name - missing")
	}
	remoteName := params[1]

	if len(params) < 3 {
		panic("remote dir - missing")
	}
	remoteDir := params[2]

	if ds.Load(dsDir) {
		remote := m.NewRemoteDirectory(remoteName, remoteDir)

		ds.AddRemote(remoteName, remote)

		err := ds.Save()

		if err != nil {
			panic(err)
		}
		ds.Save()
	}
}

func commandListRemotes(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}

	if ds.Load(params[0]) {
		for _, r := range ds.ListRemotes() {
			if r.LastSync == 0 {
				fmt.Printf("%s Last sync: never ", r.Name)
			} else {
				fmt.Printf("%s Last sync: %s ", r.Name, time.Unix(r.LastSync, 0))
			}
			if r.File != nil {
				fmt.Printf("directory: %s\n", r.File.Path)
			} else {
				fmt.Printf("remote type: %s\n", r.Target)
			}
		}
	}
}

func commandDelRemote(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}
	dsDir := params[0]

	if len(params) < 2 {
		panic("remote name - missing")
	}
	remoteName := params[1]

	if ds.Load(dsDir) {
		ds.DelRemote(remoteName)

		err := ds.Save()

		if err != nil {
			panic(err)
		}
		ds.Save()
	}
}
