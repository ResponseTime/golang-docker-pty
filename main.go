package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
			buffer := make([]byte, 32)
			c.Write([]byte("Enter the amount of time you need the container for in minutes\n"))
			n, err := c.Read(buffer)
			if err != nil {
				fmt.Println("Read error:", err)
				c.Close()
				return
			}
			timeReq, err := strconv.Atoi(strings.TrimSpace(string(buffer[:n])))
			if err != nil {
				c.Write([]byte("Time cannot be parsed"))
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeReq)*time.Minute)
			cmd := exec.Command("bash")
			ptmx, err := pty.Start(cmd)
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				<-ctx.Done()
				cancel()
				c.Write([]byte("Session ended goodbye"))
				cmd.Process.Kill()
				c.Close()
			}()
			go io.Copy(c, ptmx)
			go io.Copy(ptmx, c)
		}(conn)
	}
}
