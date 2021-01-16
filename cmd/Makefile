packages = . ../console ../connectors ../connectors/garminconnect ../connectors/directory ../models ../internal ../metadataExtractor

# First target will be running make
build: test
	go build -o fitmanager

test: format
	go test -cover $(packages)

format: clean
	go fmt $(packages)
clean:
	rm -f ../models/testdata/test*
	rm -f fitmanager

