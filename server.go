package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			serverLog("", "Added new connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				serverLog("", "A connection has terminated")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {

				select {
				case connection.data <- message:
				default:
					log.Println(string(message))
					log.Println("closing socket")
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		data := make([]byte, 4096)
		length, err := client.socket.Read(data)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			log.Println("closed connection")
			break
		}
		if length > 0 {
			var c Command
			log.Println(string(data))
			dec := json.NewDecoder(strings.NewReader(string(data)))
			err := dec.Decode(&c)
			if err != nil {
				log.Println("error while receiving data: ", err)
			}
			log.Println("Command:", c.Command)
			log.Println("Payload:", c.Payload)
			switch c.Command {
			case commandCodeLogin:
				client.handle = c.Payload
				client.data <- []byte("Succesfully logged in as: " + client.handle)
				serverLog("", "user succesfully logged in '"+client.handle+"'")
				log.Println(client)
				for connection := range manager.clients {
					if client != connection{
						connection.data <- serverResponse("", "user logged in '"+client.handle+"'")
					}
				}
			case commandCodeLogout:
				manager.unregister <- client
				client.socket.Close()
				serverLog("", "user logged out '"+client.handle+"'", )
				manager.broadcast <- serverResponse("", "server: user logged out '"+client.handle+"'")
			case commandCodeSay:
				manager.broadcast <- serverResponse(client.handle, c.Payload)
			}
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

func serverResponse(handle string, message string) []byte {
	if handle == "" {
		handle = "server"
	}
	return []byte(fmt.Sprint(time.Now().UTC().String()+ " - " + handle+ ": "+ message))
}

func serverLog(handle string, message string) {
	if handle == "" {
		handle = "server"
	}
	log.Println( handle+ ": " + message)
}

func startServerMode() {
	serverLog("", "Starting server on "+genAddress(""))
	listener, err := net.Listen("tcp", genAddress(""))
	if err != nil {
		log.Println(err)
	}
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, _ := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}
