package garminconnect

import (
	"fmt"
	"io"
	"strconv"
	"syscall"

	connect "github.com/abrander/garmin-connect"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	exportFormat = "FIT"
)

// Passing the known externalIds (File or Id)
// the function passes the unknown (new/missing) ids
// (File or Id)[]
// TODO testing the iteration
func (gci *GarminConnectConnector) FetchDiff(knownExternalIds []string) []string {
	gci.authenticate()

	found := []string{}
	i, page := 0, 50

	for more := true; more; i += page {
		activities := fetchActivities(gci.fetcher, i, page)

	outter:
		for _, a := range activities {
			id := strconv.Itoa(a.ID)
			for _, k := range knownExternalIds {
				if k == id {
					continue outter
				}
			}
			found = append(found, id)
		}

		if len(activities) == 0 {
			more = false
		}
	}
	fmt.Printf("Found %d new activities\n", len(found))
	return found
}

func (gci *GarminConnectConnector) authenticate() {
	if !gci.authenticated {
		var email string
		fmt.Print("Email: ")
		fmt.Scanln(&email)

		fmt.Print("Password: ")

		password, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			panic(err)
		}

		gci.authenticator.SetOptions(connect.Credentials(email, string(password)))
		err = gci.authenticator.Authenticate()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nSuccess\n")
	}
}

func (gci GarminConnectConnector) ExportActivity(id int, w io.Writer, f connect.ActivityFormat) error {
	fmt.Printf("Exporting activity %d :", id)
	fileName := gci.exporter.ExportActivity(id, w, f)
	fmt.Println(fileName)
	return fileName
}

func (gci *GarminConnectConnector) Import(targetDir string, externalId string) (string, error) {
	ix, _ := strconv.Atoi(externalId)
	fileName := downloadActivity(gci.exporter, gci.io, targetDir, ix)

	return fileName, nil
}
