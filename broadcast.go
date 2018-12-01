package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"

	"github.com/blackvirus18/gochat/api"
)

//Edit: Local Network broadcast address
const broadcastAddress = "192.168.1.255:33333"

// Broadcast Listener , Listens on 33333 and updates the Global Users list
func listenAndRegisterUsers() {

	// startServer to port 33333
	udpAddress, _ := net.ResolveUDPAddr("udp4", broadcastAddress)
	udpConn, err := net.ListenUDP("udp", udpAddress)
	defer udpConn.Close()
	if err != nil {
		log.Print(err)
	}

	var user api.Handle
	for {
		// read the data and add to users.
		inputBytes := make([]byte, 4096)
		length, _, _ := udpConn.ReadFromUDP(inputBytes)
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&user)

		// Ignore the user with same host
		if user.Host != MyHandle.Host {
			users.Insert(user)
		}
	}
}

// broadcastOwnHandle go-routine that publishes it's Handle on 33333
func broadcastOwnHandle() {

	// Broadcast immediately at the start
	broadcastIsAlive()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			broadcastIsAlive()
		}
	}
}

// broadcast on 33333 every 30 seconds with MyHandle(own) Handler
func broadcastIsAlive() {
	conn, err := net.Dial("udp", broadcastAddress)
	defer conn.Close()
	if err != nil {
		log.Print(err)
		return
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(MyHandle)
	conn.Write(buffer.Bytes())
	buffer.Reset()
}
