package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blackvirus18/gochat/api"
	"google.golang.org/grpc"
)

func sendChat(receiverHandle api.Handle, message string) {
	receiverConnStr := fmt.Sprintf("%s:%d", receiverHandle.Host, receiverHandle.Port)

	receiverConn, err := grpc.Dial(receiverConnStr, grpc.WithInsecure())
	defer receiverConn.Close()
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return
	}

	chatClient := api.NewGoChatClient(receiverConn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := api.ChatRequest{
		To:      &receiverHandle,
		From:    &MyHandle,
		Message: message,
	}

	_, err = chatClient.Chat(ctx, &req)
	if err != nil {
		log.Printf("ERROR: Chat(): %v", err)
		USERS.Delete(receiverHandle.Name)
	}
	return
}
