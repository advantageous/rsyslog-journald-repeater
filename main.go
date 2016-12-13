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
	"net"
)

func main() {

	host := flag.String("host", "0.0.0.0", "hostname to listen on")
	debug := flag.Bool("debug", false, "debug flag")


	defaultPort := 514

	if runtime.GOOS=="darwin" {
		defaultPort = 5514
	}
	port := flag.Int("port", defaultPort, "port to listen on")
	bufferSize := flag.Int("buffer-size", 1024, "size of buffered channel")
	listenSyslog := flag.Bool("syslog", true, "listen to syslog")
	listenJson := flag.Bool("json", true, "listen to JSON Logstash")
	jsonPort := flag.Int("json-port", 5515, "port to listen on for JSON logstash")
	jsonBufferSize := flag.Int("json-buffer-size", 1024, "max size of JSON message channel")

	flag.Parse()

	if *listenSyslog {
		if *listenJson {
			go runSyslogHandler(*host, *port, *bufferSize, *debug)
		} else {
			runSyslogHandler(*host, *port, *bufferSize, *debug)
		}
	}


	if *listenJson {

		runJsonHandler(*host, *jsonPort, *bufferSize, *debug, *jsonBufferSize)
	}

}


func runJsonHandler(host string, port, bufferSize int, debug bool, jsonBufferSize int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp",  addr)
	if err != nil {
		panic(err)
	}

	for {
		var buf = make([]byte, jsonBufferSize)
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		if !debug {
			fmt.Printf("%s\n\n", string(buf[:n]));
		}
	}
}

func runSyslogHandler(host string, port, bufferSize int, debug bool) {

	channel := make(syslog.LogPartsChannel, bufferSize)
	handler := syslog.NewChannelHandler(channel)


	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	server.ListenUDP(fmt.Sprintf("%s:%d", host, port))
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

			tag, ok := logParts["tag"].(string)
			if ok {
				if tag != "" {
					parts["SYSLOG_IDENTIFIER"] = tag
				}
			}
			journal.Send(message, journal.Priority(priority), parts)

			if debug {
				for k, v := range logParts {
					fmt.Printf("K=%s, V=%v, \t\t Type=%T \n", k, v, v)
				}
			}

		}
	}(channel)

	server.Wait()
}
