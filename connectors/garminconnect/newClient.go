package garminconnect

import (
	connect "github.com/abrander/garmin-connect"
)

func newClient() *connect.Client {
	client := connect.NewClient(
		connect.AutoRenewSession(true),
		//connect.Credentials(email, password),
	)

	// client.Authenticate()

	return client
}
