/*
Listens to incoming syslog messages and replays these messages on journal unix domain socket.

./rsyslog-journald-repeater -h
Usage of ./rsyslog-journald-repeater:
  -debug
        debug flag
  -host string
        hostname to listen on (default "0.0.0.0")
  -port int
        port to listen on (default 5514 on darwin, default 514 on Linux, Unix, etc.)
 */
package main

import (
	"gopkg.in/mcuadros/go-syslog.v2"
	"fmt"
	"github.com/coreos/go-systemd/journal"
	"time"
	"strconv"
	"flag"
	"runtime"
)

func main() {

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	defaultPort := 514

	if runtime.GOOS=="darwin" {
		defaultPort = 5514
	}

	host := flag.String("host", "0.0.0.0", "hostname to listen on")
	port := flag.Int("port", defaultPort, "port to listen on")
	debug := flag.Bool("debug", false, "debug flag")
	flag.Parse()

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	server.ListenUDP(fmt.Sprintf("%s:%d", *host, *port))
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {

			parts := make(map[string]string)

			priority, ok := logParts["severity"].(int)
			if !ok {
				priority = int(journal.PriErr)
			}

			facility, ok := logParts["facility"].(int)
			if !ok {
				facility = 6
			}
			parts["SYSLOG_FACILITY"] = strconv.Itoa(facility)

			message, ok := logParts["content"].(string)
			if !ok {
				message = "no message"
			}

			timestamp, ok := logParts["timestamp"].(time.Time)
			if !ok {
				parts["_SOURCE_REALTIME_TIMESTAMP"] = strconv.FormatInt((time.Now().UnixNano() / 1000), 10)
			} else {

				parts["_SOURCE_REALTIME_TIMESTAMP"] = strconv.FormatInt(timestamp.UnixNano() / 1000, 10)
			}

			parts["_TRANSPORT"] = "syslog"
			if !ok {
				parts["_SOURCE_REALTIME_TIMESTAMP"] = strconv.FormatInt((time.Now().UnixNano() / 1000), 10)
			} else {

				parts["_SOURCE_REALTIME_TIMESTAMP"] = strconv.FormatInt(timestamp.UnixNano() / 1000, 10)
			}

			tag, ok := logParts["tag"].(string)
			if ok {
				if tag != "" {
					parts["SYSLOG_IDENTIFIER"] = tag
				}
			}

			journal.Send(message, journal.Priority(priority), parts)

			if *debug {
				for k, v := range logParts {

					fmt.Printf("K=%s, V=%v, \t\t Type=%T \n", k, v, v)
				}
			}

		}
	}(channel)

	server.Wait()

}
