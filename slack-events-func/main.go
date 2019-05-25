package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"net/http"
	"os"
	"strings"
)

var token = os.Getenv("TOKEN")
var verifyToken = os.Getenv("VERIFY_TOKEN")

var api = slack.New(token)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// retry の場合は無視
	retryReason := request.Headers["X-Slack-Retry-Reason"]
	if retryReason != "" {
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: http.StatusOK,
		}, nil
	}

	body := request.Body
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verifyToken}))
	if e != nil {
		fmt.Println("[ERROR] ParseEvent error:" +  e.Error())
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			fmt.Println("[ERROR] ParseEvent error:" +  e.Error())
			return events.APIGatewayProxyResponse{
				Body:       "",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return events.APIGatewayProxyResponse{
			Body:      r.Challenge,
			StatusCode: http.StatusOK,
		}, nil
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			handleMention(*ev)
			return events.APIGatewayProxyResponse{
				Body:       "",
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: http.StatusOK,
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}



func handleMention(ev slackevents.AppMentionEvent) {
	text := ev.Text

	fmt.Println("mention words:" + text)

	words := strings.Fields(text)

	firstWord := words[1]

	if firstWord == "help" {
		handleHelp(ev)
		return
	}

	if firstWord == "iap" {
		handleIap(ev)
		return
	}

	if firstWord == "inhouse" {
		handleInhouse(ev)
		return
	}

	sendMsg("なにいってるかわかりません。。。",ev.Channel)
}


func handleHelp(ev slackevents.AppMentionEvent) {

	actions :=  []slack.AttachmentAction{
		{
			Name: "YES",
			Text: "はい",
			Type: "button",
			Style: "primary",
			Value: "true",
		},
		{
			Name: "NO",
			Text: "いいえ",
			Type: "button",
			Value: "false",
		},
	}

	attachment :=  slack.Attachment{
		Title: "確認",
		Text: "バージョンbuild本当ですか?",
		Actions: actions,
		Color: "#3AA3E3",
		CallbackID: "callback_help",
	}

	_, _, _, _ = api.SendMessage(ev.Channel, slack.MsgOptionAttachments(attachment))

	//sendMsg("```コマンドサンプル: \n inhouse build: @ios_slack_bot inhouse v4.8.0 \n 課金build: @ios_slack_bot iap v4.8.0 ``",ev.Channel)
}


func handleIap(ev slackevents.AppMentionEvent) {
	version := strings.Fields(ev.Text)[2]

	BuildIap(version)

	sendMsg("deploying iap build, version: " + version,ev.Channel)
}

func handleInhouse(ev slackevents.AppMentionEvent) {
	version := strings.Fields(ev.Text)[2]

	BuildInhouse(version)

	sendMsg("deploying inhouse build, version: " + version,ev.Channel)
}


func sendMsg(mes string, channel string) {
	_, _, _ = api.PostMessage(channel, slack.MsgOptionText(mes, false))
}