package metadataExtractor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/tormoder/fit"
)

func Extract(file string) (Activity string,
	Device string,
	Start int64,
	End int64,
	Samples int64,
	Bands []string,
	Created int64,
) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("%s: %s\n", file, err.Error())
		return "Error", "can not read", 0, 0, 0, []string{}, time.Now().Unix()
	}
	fitFile, err := fit.Decode(bytes.NewReader(f))
	if err != nil {
		fmt.Printf("%s: %s\n", file, err.Error())
		return "Error", "can not decode", 0, 0, 0, []string{}, time.Now().Unix()
	}
	activity, err := fitFile.Activity()
	if err != nil {
		fmt.Printf("%s: %s\n", file, err.Error())
		return "Error", "no activity", 0, 0, 0, []string{}, time.Now().Unix()
	}
	if len(activity.Sessions) == 0 {
		fmt.Printf("%s: No session found\n", file)
		return "Error", "no activity", 0, 0, 0, []string{}, time.Now().Unix()
	}
	if len(activity.Records) == 0 {
		fmt.Printf("%s: Empty activity\n", file)
		return "Error", "empty activity", 0, 0, 0, []string{}, time.Now().Unix()
	}

	session := activity.Sessions[0]

	sport := fmt.Sprint(session.Sport)
	device := fmt.Sprint(fitFile.FileId.GetProduct())
	created := fitFile.FileId.TimeCreated.Unix()
	samples := len(activity.Records)
	// calc start, end, record count
	start := activity.Records[0].Timestamp.Unix()
	end := activity.Records[samples-1].Timestamp.Unix()
	bands := []string{}
	r := activity.Records[samples-1]

	if !r.PositionLat.Invalid() && !r.PositionLong.Invalid() {
		bands = append(bands, "pos")
	}
	if session.GetAvgAltitudeScaled() == math.NaN() {
		bands = append(bands, "alt")
	}
	if session.AvgHeartRate != 0xFF && session.AvgHeartRate != 0 {
		bands = append(bands, "hr")
	}
	if session.AvgCadence != 0xFF && session.AvgCadence != 0 {
		bands = append(bands, "cad")
	}
	if session.TotalDistance != 0xFFFFFFFF && session.TotalDistance != 0 {
		bands = append(bands, "dist")
	}
	if session.AvgPower != 0xFFFF && session.AvgPower != 0 {
		bands = append(bands, "pow")
	}
	if session.TotalCalories != 0xFFFF && session.TotalCalories != 0 {
		bands = append(bands, "cal")
	}
	if session.AvgSpeed != 0xFFFF && session.AvgSpeed != 0 {
		bands = append(bands, "speed")
	}
	if session.AvgGrade != 0x7FFF && session.AvgGrade != 0 {
		bands = append(bands, "grad")
	}
	if session.AvgTemperature != 0x7F && session.AvgTemperature != 0 {
		bands = append(bands, "temp")
	}

	return sport, device, start, end, int64(samples), bands, created
}
