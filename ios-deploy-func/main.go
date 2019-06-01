package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"os"
	"slack-bot-lambda/bitrise"
)

var token = os.Getenv("TOKEN")

var api = slack.New(token)

type Event struct {
	Version string    `json:"version"`
	Name string `json:"name"`
	Type string `json:"type"`
	Channel string `json:"channel"`
	ResponseUrl string `json:"responseUrl"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func HandleRequest(request Event) (Response, error) {
	if request.Type == "iap" {
		handleIap(request)
		updateSlackMessage(request)
	}

	if request.Type == "inhouse" {
		handleInhouse(request)
		updateSlackMessage(request)
	}

	return Response{
		Message: "",
		Ok:      true,
	}, nil
}


func updateSlackMessage(event Event) {
	if event.ResponseUrl == "" {
		return
	}

	jsonStr := `{"text":"` + "Deploying " + event.Version +" "+ event.Type +" build" + `"}`

	req, _ := http.NewRequest(
		"POST",
		event.ResponseUrl,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)

	b, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		fmt.Println("resend response url :" + string(b))
	}
	defer resp.Body.Close()
}

func main() {
	lambda.Start(HandleRequest)
}

func handleIap(request Event) {
	version := request.Version

	bitrise.BuildIap(version)

	sendMsg("Deploying " + version +" iap build" ,request.Channel)
}

func handleInhouse(request Event) {
	version := request.Version

	bitrise.BuildInhouse(version)

	sendMsg("Deploying " + version +" inhouse build" ,request.Channel)
}


func sendMsg(mes string, channel string) {
	_, _, _ = api.PostMessage(channel, slack.MsgOptionText(mes, false))
}