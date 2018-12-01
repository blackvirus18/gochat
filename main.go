package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/blackvirus18/gochat/api"
)

const helpStr = `Commands
1. /users :- Get list of live users
2. @{user} message :- send message to specified user
3. /exit :- Exit the Chat
4. /all :- Send message to all the users [TODO]`

var (
	name      = flag.String("name", "", "The name you want to chat as")
	port      = flag.Int("port", 12345, "Port that your server will run on.")
	host      = flag.String("host", "", "Host IP that your server is running on.")
	stdReader = bufio.NewReader(os.Stdin)
)

//MyHandle is the description of the endpoint through which I will be able to chat with others
var MyHandle api.Handle

//users are a list of users to which I can talk on chat
var users = &PeerHandleMapSync{
	PeerHandleMap: make(map[string]api.Handle),
}

func main() {
	// Parse flags for host, port and name
	flag.Parse()

	// TODO-WORKSHOP-STEP-1: If the name and host are empty, return an error with help message

	// TODO-WORKSHOP-STEP-2: Initialize global MyHandle of type api.Handle

	// Broadcast for is-alive on 33333 with own UserHandle.
	go broadcastOwnHandle()

	// Listener for is-alive broadcasts from other hosts. Listening on 33333
	go listenAndRegisterUsers()

	// gRPC listener
	go startServer()

	go func() {
		for {
			fmt.Printf("> ")
			textInput, _ := stdReader.ReadString('\n')
			// convert CRLF to LF
			textInput = strings.Replace(textInput, "\n", "", -1)
			parseAndExecInput(textInput)
		}
	}()

	waitForExit()
}

// Handle the input chat messages as well as help commands
func parseAndExecInput(input string) {
	// Split the line into 2 tokens (cmd and message)
	tokens := strings.SplitN(input, " ", 2)
	cmd := tokens[0]

	switch {
	case cmd == "":
	case strings.ToLower(cmd) == "/users":
		fmt.Println(users)
	case strings.ToLower(cmd) == "/exit":
		os.Exit(1)
	case cmd[0] == '@':
		// TODO-WORKSHOP-STEP-9: Write code to sendChat. Example
		// "@gautam hello golang" should send a message to handle with name "gautam" and message "hello golang"
		// Invoke sendChat to send the  message
	default:
		fmt.Println(helpStr)
	}
}

func waitForExit() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...")
		wg.Done()
	}()
	wg.Wait()
	os.Exit(0)
}
