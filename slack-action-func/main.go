package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	l "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Event struct {
	Version string    `json:"version"`
	Name string `json:"name"`
	Type string `json:"type"`
	Channel string `json:"channel"`
	ResponseUrl string `json:"responseUrl"`
}

func parseBody(body string) string {
	decodedValue, _ := url.QueryUnescape(body)
	log.Println("decoded body :" + decodedValue)

	data := strings.Trim(decodedValue,":payload=")
	return data
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	log.Println("request body :" + body)

	data := parseBody(body)

	var message slack.InteractionCallback

	if err := json.Unmarshal([]byte(data), &message); err != nil {
		log.Print("json unmarshal message failed: ", err)
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if message.Type == "interactive_message" {
		log.Println("receive interactive message, name:" + message.Name + "callback id:" + message.CallbackID)

		if message.CallbackID == "callback_iap_deploy" {
			text := handleBuild(message, "iap")
			return events.APIGatewayProxyResponse{
				Body:       text,
				StatusCode: http.StatusOK,
			}, nil
		}

		if message.CallbackID == "callback_inhouse_deploy" {
			text := handleBuild(message, "inhouse")
			return events.APIGatewayProxyResponse{
				Body:       text,
				StatusCode: http.StatusOK,
			}, nil
		}


		if message.CallbackID == "callback_deploy" {
			text := handleDeploy(message)
			return events.APIGatewayProxyResponse{
				Body:       text,
				StatusCode: http.StatusOK,
			}, nil
		}



	}

	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: http.StatusOK,
	}, nil

}


func handleDeploy(message slack.InteractionCallback) string {
	action := message.ActionCallback.AttachmentActions[0]

	if action.Name == "cancel" {
		return fmt.Sprintf(":x: @%s canceled the deploy", message.User.Name)
	}

	if action.Name == "version_list" {

		value := action.SelectedOptions[0].Value

		actions :=  []slack.AttachmentAction{
			{
				Name:  "inhouse",
				Text:  "inhouse",
				Type:  "button",
				Value: value,
			},
			{
				Name:  "iap",
				Text:  "課金",
				Type:  "button",
				Value: value,
			},
		}

		originalMessage := message.OriginalMessage
		originalMessage.ResponseType = "ephemeral"
		originalMessage.Attachments[0].Text = fmt.Sprintf("%s buildするtype選択", strings.Title(value))
		originalMessage.Attachments[0].Actions = actions

		result, _ := json.Marshal(&originalMessage)
		return string(result)
	}


	if action.Name == "inhouse" {
		version := action.Value

		//bitrise.BuildInhouse(version)

		return "Deploying " + version +" inhouse build"
	}

	if action.Name == "iap" {
		version := action.Value

		//bitrise.BuildIap(version)

		return "Deploying " + version +" iap build"
	}


	return "test"

}

func responseMessage(original slack.Message, title, value string) string {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	result, _ := json.Marshal(&original)
	return string(result)
}

func handleBuild(message slack.InteractionCallback, buildType string) string {
	action := message.ActionCallback.AttachmentActions[0]
	version := action.Value

	if action.Name == "YES" {
		callDeployLambdaFunc(message, buildType)
		return ""
	}

	return "Deploying " + version + " "+ buildType +" build has been canceled"
}

func callDeployLambdaFunc(message slack.InteractionCallback, buildType string) {
	action := message.ActionCallback.AttachmentActions[0]
	version := action.Value

	event := Event{
		Version: version,
		Name: "deploy",
		Type: buildType,
		Channel: message.Channel.Name,
		ResponseUrl: message.ResponseURL,
	}

	jsonBytes, _ := json.Marshal(&event)
	svc := lambda.New(session.New())

	print("event data:" + string(jsonBytes))
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("arn:aws:lambda:ap-northeast-1:258948870772:function:ios-deploy-func"),
		Payload:        jsonBytes,
		InvocationType: aws.String("Event"),
	}

	resp, _ := svc.Invoke(input)
	fmt.Println("deploy result:" + resp.GoString())
}

func main() {
	l.Start(HandleRequest)
}