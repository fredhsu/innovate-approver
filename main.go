package main

import (
	"fmt"

	"encoding/json"

	"github.com/nats-io/nats.go"
)

type ApprovalRequest struct {
	Id       int
	Approved bool
	Request  string
}

func main() {
	currentApprovalID := 1
	nc, _ := nats.Connect("nats://localhost:4222")
	// TODO should track approval requests in an array, or later a DB possibly NATS?
	fmt.Println("Subscribing to slackbot.command")
	approvals := make(map[int]ApprovalRequest)
	nc.Subscribe("slackbot.command", func(m *nats.Msg) {
		fmt.Printf("Recevied request %s\n", string(m.Data))
		approvalRequest := ApprovalRequest{
			Id:       currentApprovalID,
			Approved: false,
			Request:  string(m.Data),
		}
		approvals[currentApprovalID] = approvalRequest
		currentApprovalID++
		fmt.Println("Publishing to slackbot.approve.request")
		// make this json
		requestmsg, _ := json.Marshal(approvalRequest)
		fmt.Println(string(requestmsg))
		nc.Publish("slackbot.approve.request", requestmsg)
	})
	fmt.Println("Subscribing to slackbot.approve.response")
	nc.Subscribe("slackbot.approve.response", func(m *nats.Msg) {
		fmt.Printf("Recevied approval %s\n", string(m.Data))
		//approvalRequest := ApprovalRequest{Id: currentApprovalID, Approved: true}
		fmt.Println("Publishing to slackbot.notify")
		nc.Publish("slackbot.notify", []byte("Request to add : has been approved"))
	})
	for {
	}
}
