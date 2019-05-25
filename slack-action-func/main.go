package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"net/url"
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
		// Order Coffeeボタンの入力を受け付けてフォームを表示する

		log.Println("receive interactive message, name:" + message.Name + "callback id:" + message.CallbackID)

		text := "了解!deploy しますよ!"
		return events.APIGatewayProxyResponse{
			Body:       text,
			StatusCode: http.StatusOK,
		}, nil
		if message.CallbackID == "callback_help" {

		}


	} else if message.Type == "dialog_submission" {
		// フォームの入力を受け付けて何かする
	}




	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}