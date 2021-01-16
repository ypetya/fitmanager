package mocks

import "log"

var FakeCalls []string

func LogCall(msg string) {
	FakeCalls = append(FakeCalls, msg)
}

func bail(err error) {
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
