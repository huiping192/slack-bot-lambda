package main

import (
	"context"
	"encoding/json"
	"fmt"
	l "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)


type Event struct {
	Version string    `json:"version"`
	Name string `json:"name"`
	Type string `json:"type"`
	Channel string `json:"channel"`
	ResponseUrl string `json:"responseUrl"`
}

type IoTButtonEvent struct {
	DeviceEvent struct {
		ButtonClicked struct {
			ClickType string `json:"clickType"`
		} `json:"buttonClicked"`
	} `json:"deviceEvent"`
}

func SendRequest(clickType string){
	var eventName string

	if clickType == "SINGLE" {
		eventName = "iap"
	} else {
		eventName = "inhouse"
	}

	callDeployLambdaFunc(eventName)
}


func callDeployLambdaFunc(buildType string) {
	version := "v4.8.0"
	channel := "jarvis-lab"

	event := Event{
		Version: version,
		Name: "deploy",
		Type: buildType,
		Channel: channel,
		ResponseUrl: "",
	}

	jsonBytes, _ := json.Marshal(&event)
	svc := lambda.New(session.New())

	print("event data:" + string(jsonBytes))
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("arn:aws:lambda:ap-northeast-1:258948870772:function:ios-deploy-func"),
		Payload:        jsonBytes,
		InvocationType: aws.String("Event"),
	}

	resp, err := svc.Invoke(input)

	if err != nil {
		fmt.Println("deploy error:" + err.Error())
	}
	fmt.Println("deploy result:" + resp.GoString())
}

func HandleRequest(ctx context.Context, event IoTButtonEvent){
	SendRequest(event.DeviceEvent.ButtonClicked.ClickType)
}

func main() {
	l.Start(HandleRequest)
}