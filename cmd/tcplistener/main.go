package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

/*

	NOTES:

	TCP:
		send data packets from one computer to another and wait for ack,
		it also wait for connection 'handshake' before the data can be sent

	UDP:
		just yeet out the data packets whether you have connection 'handshake' or not
		it also does not wait for ack. and you have to implement all the sending strategies,
		how to reconstruct it and how to handle it if a packet is missing etc.
		But it is more performant than TCP since it does not need to wait for ack.

*/

func getLinesChannel(f io.ReadCloser) <-chan string {
	c := make(chan string)
	go func() {
		defer f.Close()
		defer close(c)

		buf := make([]byte, 8)
		var msg string = ""
		for {
			n, err := f.Read(buf)
			if n > 0 {
				content := string(buf[:n])
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if i < len(lines)-1 {
						msg += line + "\n"
						c <- msg
						msg = "" // Reset msg for the next iteration
					} else {
						msg += line
					}
				}
			}
			if err != nil {
				if err == io.EOF {
					if msg != "" {
						c <- msg // Send any remaining message
					}
					return
				}
				fmt.Printf("Error: %s\n", err)
				return
			}
		}
	}()
	return c
}

func main() {

	tcp, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer tcp.Close()

	for {
		conn, err := tcp.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		defer conn.Close()

		fmt.Println("Tcp has been accepted")
		linesChannel := getLinesChannel(conn)
		for line := range linesChannel {
			fmt.Println(line)
		}
		fmt.Println("Line Channel has been closed")
	}

}
