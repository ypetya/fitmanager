package main

import m "../models"

func callFilter(ds m.IDataStore, params []string) []m.Exercise {
	f := ds.NewRemoteFilter()
	remote := "local"
	cond := []string{}
	for i, arg := range params {
		if arg[0] == '+' || arg[0] == '-' {
			cond = append(cond, params[i])
		} else if remote != "local" {
			f.Remote(remote, cond)
			remote = params[i]
			cond = []string{}
		} else {
			remote = params[i]
		}
	}
	f.Remote(remote, cond)

	return ds.Filter(f)
}
