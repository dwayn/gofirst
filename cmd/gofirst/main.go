package main

import (
	"fmt"
	"net"
	"os"

	"github.com/akamensky/argparse"
	"github.com/dwayn/gofirst/command"
	"github.com/dwayn/gofirst/protocol/barbershop"
	"github.com/dwayn/gofirst/protocol/resp"
	"github.com/dwayn/gofirst/queue"
	"github.com/dwayn/gofirst/stats"
)

func main() {
	connType := "tcp4"
	parser := argparse.NewParser("gofirst", "In memory queueing daemon, currently supports: \n\t\t\t- priority queue with in queue reprioritization")
	host := parser.String("l", "listen", &argparse.Options{Default: "localhost", Help: "Interface to listen on"})
	port := parser.Int("p", "port", &argparse.Options{Default: 3333, Help: "Port to listen on"})
	protocol := parser.Selector("r", "protocol", []string{"resp", "barbershop"}, &argparse.Options{Default: "resp", Help: "Network protocol: resp, barbershop"})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	var connectionHandler func(net.Conn, chan command.Request, chan stats.Metric)

	l, err := net.Listen(connType, fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	switch *protocol {
	case "resp":
		connectionHandler = resp.HandleConnection
	case "barbershop":
		connectionHandler = barbershop.HandleConnection
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
