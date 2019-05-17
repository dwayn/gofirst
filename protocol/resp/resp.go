package resp

// Fully RESP compatible (https://redis.io/topics/protocol) protocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/dwayn/gofirst/command"
	"github.com/dwayn/gofirst/stats"
)

// HandleConnection handles a freshly connected client and manages all protocol interactions with the client
func HandleConnection(c net.Conn, queueManager chan command.Request, metrics chan stats.Metric) {
	metrics <- stats.Metric{Metric: "connected", Op: stats.Add, Value: 1}
	metrics <- stats.Metric{Metric: "connections", Op: stats.Add, Value: 1}
	// using a buffered channel so that the central queue manager thread does not get
	// 	blocked waiting on a connection handler to take a message off of the channel
	responseChannel := make(chan command.Response, 1)
	metricsResponse := stats.CreateReplyChannel()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err == io.EOF {
			fmt.Println("Client disconnected:", c.RemoteAddr().String())
			metrics <- stats.Metric{Metric: "connected", Op: stats.Sub, Value: 1}
			break
		}
		if err != nil {
			fmt.Println("Error reading from socket", err)
			metrics <- stats.Metric{Metric: "connected", Op: stats.Sub, Value: 1}
			break
		}
		metrics <- stats.Metric{Metric: "ops", Op: stats.Add, Value: 1}
		temp := strings.TrimSpace(string(netData))
		// this should contain the reply that will be written to the socket
		var replyString string
		defaultResponseBody := ""
		skipCommonHandler := false
		responseTypePrefix := "+"
		switch {
		case strings.HasPrefix(strings.ToLower(temp), "update "):
			// format: UPDATE key priority
			request := command.Request{OpType: command.Update, ResponseChannel: responseChannel}
			parts := strings.Split(temp, " ")
			commandOk := false
			switch len(parts) {
			case 3:
				v, err := strconv.Atoi(parts[2])
				if err == nil {
					commandOk = true
					request.ItemValue = v
					request.ItemKey = parts[1]
				}
			case 2:
				commandOk = true
				request.ItemValue = 1
				request.ItemKey = parts[1]
			}
			if commandOk {
				queueManager <- request
				responseTypePrefix = ":"
			} else {
				replyString = "-999 ERROR Invalid command format\r\n"
				skipCommonHandler = true
			}

		case strings.ToLower(temp) == "peek":
			// format: PEEK
			request := command.Request{OpType: command.Peek, ResponseChannel: responseChannel}
			queueManager <- request
		case strings.HasPrefix(strings.ToLower(temp), "score "):
			// format: SCORE key
			request := command.Request{OpType: command.Score, ResponseChannel: responseChannel}
			commandOk := false
			parts := strings.Split(temp, " ")
			if len(parts) == 2 {
				commandOk = true
				request.ItemKey = parts[1]
				responseTypePrefix = ":"
			}
			if commandOk {
				queueManager <- request
			} else {
				replyString = "-999 ERROR Invalid command format\r\n"
				skipCommonHandler = true
			}
		case strings.ToLower(temp) == "next":
			// format: NEXT
			request := command.Request{OpType: command.Next, ResponseChannel: responseChannel}
			queueManager <- request
		case strings.ToLower(temp) == "info":
			// format: INFO
			replyString = ""
			metrics <- stats.Metric{Metric: "*", Op: stats.Get, Resp: metricsResponse}
			metricCount := 0
			for {
				m := <-metricsResponse
				if m.Metric == "" {
					break
				}
				metricCount++
				replyString = fmt.Sprintf("%s$%d\r\n%s\r\n:%d\r\n", replyString, len(m.Metric), m.Metric, m.Value)
			}
			replyString = fmt.Sprintf("*%d\r\n%s", metricCount, replyString)
			skipCommonHandler = true

		case strings.ToLower(temp) == "close":
			metrics <- stats.Metric{Metric: "connected", Op: stats.Sub, Value: 1}
			return
		default:
			replyString = "-999 ERROR Invalid command format\r\n"
			// skipping the common handler her because there is an error figuring out what command requested is
			skipCommonHandler = true
		}
		if !skipCommonHandler {
			// waits until receives a response from queue thread
			reply := <-responseChannel

			switch reply.ErrorCode {
			case 0:

				if reply.ResponseBody != "" {
					replyString = fmt.Sprintf("%s%s\r\n", responseTypePrefix, reply.ResponseBody)
				} else {
					replyString = fmt.Sprintf("%s%s\r\n", responseTypePrefix, defaultResponseBody)
				}
			case command.ErrQueueEmpty:
				replyString = "$-1\r\n"
			case command.ErrNotFound:
				replyString = "$-1\r\n"
			default:
				replyString = fmt.Sprintf("-%d %s\r\n", reply.ErrorCode, reply.ErrorMessage)
			}
		}

		c.Write([]byte(replyString))

	}
}
