package main

import (
	"fmt"
	"strconv"
	"time"

	m "github.com/ypetya/fitmanager/models"
)

func commandList(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}

	if ds.Load(params[0]) {
		filtered := []m.Exercise{}
		if len(params) > 1 {
			filtered = callFilter(ds, params[1:])
		} else {
			filtered = ds.GetExercises()
		}
		for _, e := range filtered {
			duration := formatDuration(e.Meta.Start, e.Meta.End)
			created := formatCreated(e.Meta.Created)
			ref, _ := e.GetLocalRef()
			fmt.Println(created, e.Meta.Activity, duration, ref, e.OverlappingIds)
		}
	}
}

func formatCreated(created int64) string {
	return time.Unix(created, 0).Format("2006-01-02 15:04")
}

func formatDuration(start int64, end int64) string {
	diff := end - start
	ret := ""

	if diff > 0 {
		mins, seconds := diff/60, diff%60
		hours, mins := mins/60, mins%60

		if hours > 0 {
			ret += strconv.Itoa(int(hours)) + "h"
		}
		if mins > 0 {
			ret += strconv.Itoa(int(mins)) + "m"
		}
		if seconds > 0 {
			ret += strconv.Itoa(int(seconds)) + "s"
		}
	} else {
		ret = "0"
	}

	return ret
}
