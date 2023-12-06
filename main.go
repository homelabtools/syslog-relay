package main

import (
	"fmt"
	"strconv"
	"time"

	//"github.com/coreos/go-systemd/v22/journal"

	"github.com/araddon/dateparse"
	"github.com/coreos/go-systemd/journal"
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
			severity, err := strconv.Atoi(fmt.Sprint(logParts["severity"]))
			if err != nil {
				severity = int(journal.PriInfo)
			}
			tstamp, err := dateparse.ParseAny(timestamp)
			if err != nil {
				tstamp = time.Now()
			}

			vars := map[string]string{
				"UNIT":              "syslog-relay",
				"SYSLOG_IDENTIFIER": tag,
				"SYSLOG_TIMESTAMP":  tstamp.Format(time.RFC3339),
			}
			sev := journal.Priority(severity)
			err = journal.Print(sev, fmt.Sprintf("client: %q, time: %q, severity: %q, message: %q", client, tstamp, sevToString(sev), content), vars)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(channel)

	server.Wait()
	return nil
}

func sevToString(sev journal.Priority) string {
	switch sev {
	case journal.PriDebug:
		return "DEBUG"
	case journal.PriInfo:
		return "INFO"
	case journal.PriNotice:
		return "NOTICE"
	case journal.PriWarning:
		return "WARNING"
	case journal.PriErr:
		return "ERROR"
	case journal.PriAlert:
		return "ALERT"
	case journal.PriCrit:
		return "CRITICAL"
	default:
		panic(fmt.Sprintf("Unknown severity: %d", sev))
	}

}

//map[client:192.168.0.1:59410 content:httpds 1621:notify_rc restart_logger facility:1 hostname:router-9599A2F-C priority:15 severity:7 tag:rc_service timestamp:2023-07-30 16:23:39 +0000 UTC tls_peer:]
//map[client:192.168.0.1:55211 content:klogd started: BusyBox v1.25.1 (2023-05-07 12:35:02 EDT) facility:1 hostname:router-9599A2F-C priority:13 severity:5 tag:kernel timestamp:2023-07-30 16:23:39 +0000 UTC tls_peer:]
//map[client:192.168.0.1:55211 content:bsd: Sending act Frame to a4:cf:99:68:72:79 with transition target eth5 ssid fc:34:97:2c:1c:b0 facility:1 hostname:router-9599A2F-C priority:8 severity:0 tag:bsd timestamp:2023-07-30 16:23:43 +0000 UTC tls_peer:]
//map[client:192.168.0.1:55211 content:bsd: BSS Transit Response: ifname=eth6, event=156, token=41, status=6, mac=34:10:fc:34:97:2c facility:1 hostname:router-9599A2F-C priority:8 severity:0 tag:bsd timestamp:2023-07-30 16:23:43 +0000 UTC tls_peer:]
//map[client:192.168.0.1:55211 content:bsd: BSS Transit Response: STA reject facility:1 hostname:router-9599A2F-C priority:8 severity:0 tag:bsd timestamp:2023-07-30 16:23:43 +0000 UTC tls_peer:]
//map[client:192.168.0.1:55211 content:bsd: Skip STA:a4:cf:99:68:72:79 reject BSSID facility:1 hostname:router-9599A2F-C priority:8 severity:0 tag:bsd timestamp:2023-07-30 16:23:43 +0000 UTC tls_peer:]
//
