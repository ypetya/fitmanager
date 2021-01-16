package garminconnect

import (
	"fmt"

	connect "github.com/abrander/garmin-connect"

	"../../internal"
)

func downloadActivity(
	client IActivityExporter,
	io internal.IFileCreator,
	dir string,
	activityID int) (fileName string) {

	format, err := connect.FormatFromExtension(exportFormat)
	if err != nil {
		fmt.Println(err.Error())
		return "invalid"
	}
	name := fmt.Sprintf("%d.%s", activityID, format.Extension())

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Unexpected in downloadActivity call", dir, activityID, r)
			fileName = name
		}
	}()

	if !io.FileExists(dir + name) {
		f := io.CreateFileForWriting(dir + name)
		err = client.ExportActivity(activityID, f, format)

		if err != nil {
			fmt.Println("ExportActivity failed", err.Error())
		}
	}

	return name
}

func fetchActivities(client IActivitiesFetcher, offset int, count int) (returnArr []connect.Activity) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Unexpected in fetchActivities call:", r)
			returnArr = []connect.Activity{}
		}
	}()
	activities, err := client.Activities("", offset, count)
	if err != nil {
		fmt.Println(err.Error())
		return []connect.Activity{}
	}

	return activities
}
