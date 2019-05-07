package main

import (
	"fmt"
	"net"

	"github.com/dwayn/gofirst/protocol/resp"
)

// TODO: move these to runtime/config options
const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "3333"
	CONN_TYPE        = "tcp4"
	PROTOCOL_HANDLER = "resp"
)

func main() {
	// arguments := os.Args
	// if len(arguments) == 1 {
	// 		fmt.Println("Please provide a port number!")
	// 		return
	// }

	// PORT := ":" + arguments[1]
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch PROTOCOL_HANDLER {
		case "resp":
			go resp.HandleConnection(c)
			// go resp.HandleConnection(c)
		}
	}
}
