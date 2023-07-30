package main

import (
	"fmt"

	"github.com/coreos/go-systemd/v22/journal"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	err := mainE()
	if err != nil {
		panic(err)
	}
}

func mainE() error {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:514")
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			client := logParts["client"]
			timestamp := logParts["timestamp"]
			content := logParts["content"]

			fmt.Printf("%s %s %s\n", client, timestamp, content)
			err := journal.Print(journal.PriInfo, fmt.Sprintf("%s %s %s", client, timestamp, content))
			if err != nil {
				fmt.Println(err)
			}
		}
	}(channel)

	server.Wait()
	return nil
}
