package resp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// HandleConnection handles a connecting client
func HandleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err == io.EOF {
			fmt.Println("Client disconnected:", c.RemoteAddr().String())
			break
		}
		if err != nil {
			fmt.Println("Error reading from socket", err)
			break
		}

		temp := strings.TrimSpace(string(netData))
		switch {
		case strings.HasPrefix(temp, "UPDATE"):
			// format: UPDATE key priority
			c.Write([]byte("received UPDATE\r\n"))
		case strings.HasPrefix(temp, "PEEK"):
			// format: PEEK
			c.Write([]byte("received PEEK\r\n"))
		case strings.HasPrefix(temp, "SCORE"):
			// format: SCORE key
			c.Write([]byte("received SCORE\r\n"))
		case strings.HasPrefix(temp, "NEXT"):
			// format: NEXT
			c.Write([]byte("received NEXT\r\n"))
		case strings.HasPrefix(temp, "INFO"):
			// format: INFO
			c.Write([]byte("received INFO\r\n"))
		case temp == "CLOSE":
			return
		default:
			c.Write([]byte("-UNKNOWN_COMMAND\r\n"))
		}

		// c.Write([]byte(string(result)))
	}
}
