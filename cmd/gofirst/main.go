package main

import (
	"fmt"
	"net"

	"github.com/dwayn/gofirst/command"
	"github.com/dwayn/gofirst/protocol/resp"
	"github.com/dwayn/gofirst/queue"
	"github.com/dwayn/gofirst/stats"
)

// TODO: move these to runtime/config options
const (
	CONN_HOST        = "0.0.0.0"
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
	var connectionHandler func(net.Conn, chan command.Request, chan stats.Metric)

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	switch PROTOCOL_HANDLER {
	case "resp":
		connectionHandler = resp.HandleConnection
	}

	commandChannel := queue.CreateChannel()
	metrics := stats.CreateMetricChannel()
	internalQueue := queue.PriorityQueue{Metrics: metrics}

	go stats.ProcessMetrics(metrics)
	go queue.RunQueue(commandChannel, &internalQueue)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go connectionHandler(c, commandChannel, metrics)
	}
}
