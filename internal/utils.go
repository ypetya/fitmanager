package internal

import "log"

func bail(err error) {
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
