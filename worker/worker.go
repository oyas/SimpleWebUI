package worker

import (
	"fmt"
	"time"
	"net"
	"os"

	"../websocket_server"
	"../int_trap"
)

const (
	writePeriod = 2 * time.Second

	socketName = "io.socket"
)

type WorkerClient struct {
	websocket_server.Client
}


func Run(hub *websocket_server.Hub){
	//ticker := time.NewTicker(writePeriod)
	//defer ticker.Stop()

	listener, err := net.Listen("unix", socketName)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	cb := int_trap.RegisterFunc(func() {
		fmt.Println("Called cleanup")
		listener.Close()
		os.Remove(socketName)
	})
	defer int_trap.Call(cb)

	newConns := make(chan net.Conn, 128)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				break
			}
			newConns <- conn
		}
	}()

	for {
		select {
		case conn := <-newConns:
			go connection(conn, hub)
		//case now := <-ticker.C:
			//message := now.Format(time.RFC3339)
			//hub.Broadcast <- []byte(message)
		}
	}
}

func connection(conn net.Conn, hub *websocket_server.Hub) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())

	client := &WorkerClient{websocket_server.Client{Hub: hub, Send: make(chan []byte, 256)}}
	client.Hub.Register <- &client.Client
	defer func() {
		client.Hub.Unregister <- &client.Client
		conn.Close()
	}()

	var recv_message string
	buf := make([]byte, 1000000)
	ch := make(chan []byte)
	eCh := make(chan error)
	go func(ch chan []byte, eCh chan error) {
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				eCh <- err
			}
			recv_message += string(buf[:n])
			if buf[n-1] == '\n' {
				fmt.Print(string(recv_message))
				ch <- []byte(recv_message)
				recv_message = ""
			}
		}
		close(ch)
		close(eCh)
	}(ch, eCh)

	// loop for receiving unix socket data
	for {
		select {
		case data, alive := <-ch:
			hub.Broadcast <-data
			message := "sent "
			conn.Write([]byte(message))
			if alive == false {
				fmt.Printf("Connection closed\n")
				return
			}
		case err := <-eCh:
			fmt.Printf("Read error: %s\n", err)
		case hub_data := <-client.Send:
			conn.Write(hub_data)
			conn.Write([]byte("\n"))
		}
	}
}
