package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(c net.Conn) {
			// for {
			// 	buffer := make([]byte, 1024)
			// 	n, err := c.Read(buffer)
			// 	if err != nil {
			// 		fmt.Println("Read error:", err)
			// 		c.Close()
			// 		return
			// 	}
			// 	fmt.Printf("Received: %s\n", string(buffer[:n]))
			// 	c.Write([]byte("hello from tcp server\n"))
			// }
			cmd := exec.Command("bash")
			ptmx, err := pty.Start(cmd)
			if err != nil {
				fmt.Println(err)
			}
			go io.Copy(c, ptmx)
			go io.Copy(ptmx, c)
		}(conn)
	}
}
