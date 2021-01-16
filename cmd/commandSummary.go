package main

import (
	"fmt"

	m "../models"
)

func commandSummary(params []string, s interface{}) {
	ds := s.(m.IDataStore)

	if len(params) < 1 {
		panic("data-store directory - missing")
	}

	if ds.Load(params[0]) {
		exercises := ds.GetExercises()
		fmt.Printf("%d exercises in db:\n", len(exercises))
		devices := []string{}
		samples := []int64{}
		exNum := []int{}
		// where end<start
		var startEndErrNum int
		var errNum int
		var overlapNum int
		var totalSamples int64
		var notOnGarminConnect int
		var uniqueNotOnGarminConnect int
	outer:
		for _, ex := range exercises {
			if ex.Meta.Activity == "Error" {
				errNum += 1
				continue
			}
			if ex.Meta.End < ex.Meta.Start {
				startEndErrNum++
			}
			overlaps := len(ex.OverlappingIds)
			overlapNum += overlaps
			if _, err := ex.GetRemoteRef("garmin"); err != nil {
				notOnGarminConnect += 1
				if overlaps == 0 {
					uniqueNotOnGarminConnect++
				}
			}
			totalSamples += ex.Meta.Samples

			for i, device := range devices {
				if device == ex.Meta.Device {
					samples[i] += ex.Meta.Samples
					exNum[i] += 1
					continue outer
				}
			}
			devices = append(devices, ex.Meta.Device)
			samples = append(samples, ex.Meta.Samples)
			exNum = append(exNum, 1)
		}
		fmt.Println("Devices:")
		for i, v := range devices {
			fmt.Printf("%s recorded %d samples in %d exercises.\n", v, samples[i], exNum[i])
		}
		fmt.Printf("%d exercises found with meta-data errors.\n", errNum)
		fmt.Printf("%d exercises found with meta-data start end time errors.\n", startEndErrNum)
		fmt.Printf("%d overlapping exercises\n", overlapNum/2)
		fmt.Println("Total samples: ", totalSamples)
		fmt.Printf("%d exercises missing garmin-connect remote. non overlapping %d\n", notOnGarminConnect, uniqueNotOnGarminConnect)
	}
}
