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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				panic(err)
			}
			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: "ubuntu",
				Cmd:   []string{"sleep", "infinity"},
			}, nil, nil, nil, "")

			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
				panic(err)
			}

			cmd := exec.Command(
				"docker", "exec", "-it",
				resp.ID,
				"bash",
			)
			ptmx, err := pty.Start(cmd)
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				<-ctx.Done()
				cancel()
				cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
				ptmx.Close()
				c.Write([]byte("Session ended goodbye\n"))
				cmd.Process.Kill()
				c.Close()
			}()
			go io.Copy(c, ptmx)
			go io.Copy(ptmx, c)
		}(conn)
	}
}
