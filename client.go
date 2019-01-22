package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type Client struct {
	socket net.Conn
	data   chan []byte
	handle string
}

type Command struct {
	Command string
	target  string
	Payload string
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			fmt.Println(err)
			log.Println(err)
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println(string(message))
			log.Println(string(message))
		}
	}
}

func sendCommand(connection net.Conn, code string, target string, payload string) {
	jsonData, _ := json.Marshal(Command{code, target, payload})
	_, err := connection.Write(jsonData)
	if err != nil {
		fmt.Println("error while sending command: ", err)
	}
}

func startClientMode(userHandle string) {
	if userHandle == "" {
		fmt.Print("Input user handle:")
		reader := bufio.NewReader(os.Stdin)
		userHandle, _ = reader.ReadString('\n')
		userHandle = strings.TrimSuffix(userHandle, "\n")
	}
	fmt.Println("Starting client on " + genAddress("localhost"))
	connection, err := net.Dial("tcp", genAddress("localhost"))
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
	client := &Client{socket: connection}
	go client.receive()
	sendCommand(connection, commandCodeLogin, "", userHandle)
ClientLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		command, target, payload, err := clientInput(message)
		log.Println(command, target, payload)
		if err == nil {
			switch command {
			case chatCommandLogout:
				sendCommand(connection, commandCodeLogout, "", "")
				client.socket.Close()
				log.Println("closed socket")
				break ClientLoop
			case chatCommandSay:
				if payload != "" {
					sendCommand(connection, commandCodeSay, "", payload)
				}
				break
			case chatCommandUserList:
				sendCommand(connection, commandCodeUserList, "", "")
				break
			}
		} else {
			fmt.Println("something went wrong with interpreting the text input:", err)
		}
	}
}

func clientInput(input string) (string, string, string, error) {
	regex := regexp.MustCompile(`^(?:\/(?P<command>\w+))?(?: +\[(?P<target>\w+)\] +)?(?: *(?P<payload>[^\n]+))?`)
	if regex.MatchString(input) {
		res := regex.FindStringSubmatch(input)
		command := res[1]
		target := res[2]
		payload := res[3]
		if command == "" || command == chatShortCommandSay {
			command = chatCommandSay
		}
		return command, target, payload, nil
	} else {
		return "", "", "", errors.New("no valid match")
	}
}
