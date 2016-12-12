package main


import (

	"net"
	"os"
	"fmt"
)

func main() {
	conn, err := net.ListenUnixgram("unixgram",  &net.UnixAddr{"/run/systemd/journal/socket", "unixgram"})
	if err != nil {
		panic(err)
	}
	defer os.Remove("/run/systemd/journal/socket")

	for {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(buf[:n]));
	}
}