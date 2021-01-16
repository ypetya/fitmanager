package main

import (
	"fmt"
	"os"

	"../connectors"
	dirConnector "../connectors/directory"
	garminConnector "../connectors/garminconnect"
	console "../console"
	m "../models"

	"github.com/ypetya/fitfixer"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed:", r)
			printUsage()
		}
	}()
	// GCIO is garmin connect importer, exporter and metadata extractor
	gcio := garminConnector.NewGarminConnectConnector()

	connectors := map[m.RemoteType]connectors.IConnector{
		m.GarminConnect: gcio,
		m.Directory:     dirConnector.NewDirectoryConnector(),
	}

	hrFixer := fitfixer.HrEnhancer{}
	hrEnhancer := m.Enhancer{
		Filter:   FilterHrEnhanceable{},
		Function: hrFixer,
		Name:     "HR",
	}

	ds := m.NewDataStore(connectors, gcio,
		[]m.Enhancer{
			hrEnhancer,
		})

	cmd := console.NewCommand(ds)

	defineCommands(&cmd)

	fmt.Printf("Fitmanager version %d\n", ds.Version)

	if !cmd.Execute(os.Args[1], os.Args[2:]) {
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`Usage:
    fitmanager <command with arguments>
    
    Command is one of:

    init <data-store directory> 
    import <data-store directory> [remote name, default: garmin]
		export <data-store directory> [remote name]
    add <data-store directory> <remote name> <directory>
    remove <data-store directory> <remote name>
    remotes <data-store directory>
    summary <data-store directory>
    list <data-store directory> [remote ... ]
    `)
	//export <data-store directory> <remote name> - Not implemented yet!
}

func defineCommands(e console.ICommands) {
	e.AddCommand("init", commandInit)
	e.AddCommand("import", commandImport)
	e.AddCommand("add", commandAddRemote)
	e.AddCommand("remove", commandDelRemote)
	e.AddCommand("remotes", commandListRemotes)
	e.AddCommand("export", commandExport)
	e.AddCommand("summary", commandSummary)
	e.AddCommand("list", commandList)
}
