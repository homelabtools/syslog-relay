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
	err := server.ListenUDP("0.0.0.0:514")
	if err != nil {
		return err
	}
	err = server.Boot()
	if err != nil {
		return err
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			client := fmt.Sprint(logParts["client"])
			timestamp := fmt.Sprint(logParts["timestamp"])
			content := fmt.Sprint(logParts["content"])
			tag := fmt.Sprint(logParts["tag"])

			vars := map[string]string{
				"UNIT":              "syslog-relay",
				"SYSLOG_IDENTIFIER": tag,
				"SYSLOG_TIMESTAMP":  timestamp,
			}

			fmt.Println(logParts)
			// TODO: extract priority
			err := journal.Print(journal.PriInfo, fmt.Sprintf("%s %s %s", client, timestamp, content), vars)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(channel)

	server.Wait()
	return nil
}
