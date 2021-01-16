package main

import (
	"fmt"

	m "github.com/ypetya/fitmanager/models"
)

func commandInit(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}

	if !ds.Load(params[0]) {
		fmt.Println("Failed to load ds.json - initializing new data-store")
		ds.AddRemote("garmin", m.AnyRemote{
			Target: m.GarminConnect,
			Name:   "garmin",
		})
	}
	ds.Save()
}
