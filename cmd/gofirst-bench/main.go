package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("gofirst-bench", "bench testing tool for gofirst queue")
	threads := parser.Int("t", "threads", &argparse.Options{Default: 10, Help: "Number of concurrent threads"})
	host := parser.String("H", "host", &argparse.Options{Default: "localhost"})
	port := parser.Int("p", "port", &argparse.Options{Default: 3333})
	items := parser.Int("i", "items", &argparse.Options{Default: 1000, Help: "Number of items to load into the queue per thread"})
	ops := parser.Int("o", "operations", &argparse.Options{Default: 0, Help: "Number of operations to do after initial queue load"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	fmt.Printf("%s:%d %d\n", *host, *port, *threads)
	var wg sync.WaitGroup
	for x := 0; x < *threads; x++ {
		wg.Add(1)
		go func(h string, p, i, o int) {
			defer wg.Done()
			doWork(h, p, i, o)
		}(*host, *port, *items, *ops)
	}
	wg.Wait()
}

func doWork(host string, port, items, ops int) {

	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	readBuffer := bufio.NewReader(conn)
	for i := 0; i < items; i++ {
		rint := rand.Int()
		fmt.Fprintf(conn, "UPDATE %d %d\r\n", rint, rint)
		readBuffer.ReadString('\n')
	}
	for i := 0; i < ops; i++ {
		if rand.Intn(1) == 0 {
			rint := rand.Int()
			fmt.Fprintf(conn, "UPDATE %d %d\r\n", rint, rint)
			readBuffer.ReadString('\n')
		} else {
			fmt.Fprintf(conn, "NEXT\r\n")
			readBuffer.ReadString('\n')
		}
	}

}
