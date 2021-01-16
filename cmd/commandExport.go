package main

import (
	"fmt"

	m "../models"
)

func commandExport(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}
	if len(params) < 2 {
		panic("remote name - missing")
	}
	remote := params[1]
	if ds.Load(params[0]) {
		fmt.Println("*** Exporting to: ", remote)
		err := ds.Export(remote)

		if err != nil {
			panic(err)
		}
		ds.Save()
	} else {
		fmt.Println("Can not load data-store from", params[0])
	}
}
