package main

import (
	"fmt"

	m "github.com/ypetya/fitmanager/models"
)

func commandImport(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}

	if ds.Load(params[0]) {
		remote := "garmin"
		if len(params) > 1 {
			remote = params[1]
		}
		fmt.Println("*** Importing from: ", remote)
		err := ds.Import(remote)

		if err != nil {
			panic(err)
		}
		ds.Save()
	} else {
		fmt.Println("Can not load data-store from", params[0])
	}
}
