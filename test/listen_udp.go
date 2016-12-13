package main


import (

	"net"
	"fmt"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:5514")

	if err != nil {
		panic(err)
	}


	conn, err := net.ListenUDP("udp",  addr)
	if err != nil {
		panic(err)
	}

	for {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n\n", string(buf[:n]));
	}
}