package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"net/url"
	"slack-bot-lambda/bitrise"
	"strings"
)

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
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if message.Type == "interactive_message" {
		log.Println("receive interactive message, name:" + message.Name + "callback id:" + message.CallbackID)

		if message.CallbackID == "callback_iap_deploy" {
			text := handleIap(message)
			return events.APIGatewayProxyResponse{
				Body:       text,
				StatusCode: http.StatusOK,
			}, nil
		}

		if message.CallbackID == "callback_inhouse_deploy" {
			text := handleInhouse(message)
			return events.APIGatewayProxyResponse{
				Body:       text,
				StatusCode: http.StatusOK,
			}, nil
		}


	} else if message.Type == "dialog_submission" {
		// フォームの入力を受け付けて何かする
	}


	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: http.StatusOK,
	}, nil

}


func handleIap(message slack.InteractionCallback) string {
	action := message.ActionCallback.AttachmentActions[0]

	version := action.Value

	if action.Name == "YES" {

		bitrise.BuildIap(version)

		return "Deploying " + version +" iap build"
	}

	return "Deploying " + version + " iap build has been canceled"
}


func handleInhouse(message slack.InteractionCallback) string {
	action := message.ActionCallback.AttachmentActions[0]

	version := action.Value

	if action.Name == "YES" {

		bitrise.BuildInhouse(version)

		return "Deploying " + version +" inhouse build"
	}

	return "Deploying " + version + " inhouse build has been canceled"
}

func main() {
	lambda.Start(HandleRequest)
}