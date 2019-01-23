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

func (client *Client) sendCommand(command *Command){
	_, err := client.socket.Write(command.encode())
	if err != nil {
		fmt.Println("error while sending command: ", err)
	}
}

type Command struct {
	Command string
	Target  string
	Payload string
}

func (command *Command) encode() []byte {
	jsonData,err := json.Marshal(*command)
	if err != nil{
		log.Println("error encoding command",err)
		log.Println(command)
	}
	return jsonData
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
	client.sendCommand(&Command{commandCodeLogin, "", userHandle})
ClientLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		command, err := clientInput(message)
		log.Println(command)
		if err == nil {
			switch command.Command {
			case commandCodeLogout:
				client.sendCommand(&command)
				break ClientLoop
			default:
				client.sendCommand(&command)
			}
		} else {
			fmt.Println("something went wrong with interpreting the text input:", err)
			break ClientLoop
		}
	}
}

func clientInput(input string) (Command, error) {
	regex := regexp.MustCompile(`^(?:\/(?P<command>\w+))?(?: +\[(?P<target>\w+)\] +)?(?: *(?P<payload>[^\n]+))?`)
	if regex.MatchString(input) {
		res := regex.FindStringSubmatch(input)
		var command string
		switch res[1] {
		case "": fallthrough
		case chatShortCommandSay: fallthrough
		case chatCommandSay:
			if res[3] == ""{
				return Command{}, errors.New("empty message")
			}
			command = commandCodeSay
		case chatCommandLogout: command = commandCodeLogout
		case chatCommandUserList: command = commandCodeUserList
		case chatShortCommandWisper: fallthrough
		case chatCommandWisper:
			if res[3] == ""{
				return Command{}, errors.New("empty message")
			}
			command = commandCodeWisper
		default:
			return Command{},errors.New("invalid Command")
		}

		return Command{command, res[2], res[3]}, nil
	} else {
		return Command{}, errors.New("invalid syntax")
	}
}
