
## rsyslog repeater. 

Listens to incoming syslog messages and replays these messages to the systemd journald.


#### Usage 
```sh
./rsyslog-journald-repeater -h
Usage of ./rsyslog-journald-repeater:
  -debug
        debug flag
  -host string
        hostname to listen on (default "0.0.0.0")
  -port int
        port to listen on (default 5514 on darwin, default 514 on Linux, Unix, etc.)
```


