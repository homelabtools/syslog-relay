package main

import (
	"fmt"

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
			fmt.Printf("%s %s %s\n", logParts["client"], logParts["timestamp"], logParts["content"])
		}
	}(channel)

	server.Wait()
	return nil
}
