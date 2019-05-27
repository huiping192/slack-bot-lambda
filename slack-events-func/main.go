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

	if firstWord == "deploy" {
		handleDeploy(ev)
		return
	}

	sendMsg("なにいってるかわかりません。。。",ev.Channel)
}


func handleHelp(ev slackevents.AppMentionEvent) {
	sendMsg("```コマンドサンプル: \n inhouse build: @ios_slack_bot inhouse v4.8.0 \n 課金build: @ios_slack_bot iap v4.8.0 ``",ev.Channel)
}

func handleDeploy(ev slackevents.AppMentionEvent) {

	versionMenu := []slack.AttachmentActionOption{
		{
			Text: "v4.8.0",
			Value: "v4.8.0",
		},
		{
			Text: "v4.7.9",
			Value: "v4.7.9",
		},
		{
			Text: "v4.7.8",
			Value: "v4.7.8",
		},
		{
			Text: "v4.7.7",
			Value: "v4.7.7",
		},
		{
			Text: "v4.7.6",
			Value: "v4.7.6",
		},
	}
	//
	//deployTypeList := []slack.AttachmentActionOption{
	//	{
	//		Text: "inhouse",
	//		Value: "inhouse",
	//	},
	//	{
	//		Text: "課金",
	//		Value: "iap",
	//	},
	//}

	actions :=  []slack.AttachmentAction{
		{
			Name: "version_list",
			Text: "バージョン選択",
			Type: "select",
			Options: versionMenu,
		},
		{
			Name:  "cancel",
			Text:  "Cancel",
			Type:  "button",
			Style: "danger",
		},
	}

	attachment :=  slack.Attachment{
		Text: "どのiOSバージョンdeployしますか?",
		Actions: actions,
		Color: "#3AA3E3",
		CallbackID: "callback_deploy",
	}

	_, _, _, _ = api.SendMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}


func handleIap(ev slackevents.AppMentionEvent) {
	version := strings.Fields(ev.Text)[2]

	actions :=  []slack.AttachmentAction{
		{
			Name: "YES",
			Text: "はい",
			Type: "button",
			Style: "primary",
			Value: version,
		},
		{
			Name: "NO",
			Text: "いいえ",
			Type: "button",
			Value: version,
		},
	}

	attachment :=  slack.Attachment{
		Title: "確認",
		Text: version+"の課金バージョン作りますが、よろしいですか?",
		Actions: actions,
		Color: "#3AA3E3",
		CallbackID: "callback_iap_deploy",
	}

	_, _, _, _ = api.SendMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}

func handleInhouse(ev slackevents.AppMentionEvent) {
	version := strings.Fields(ev.Text)[2]

	actions :=  []slack.AttachmentAction{
		{
			Name: "YES",
			Text: "はい",
			Type: "button",
			Style: "primary",
			Value: version,
		},
		{
			Name: "NO",
			Text: "いいえ",
			Type: "button",
			Value: version,
		},
	}

	attachment :=  slack.Attachment{
		Title: "確認",
		Text: version+"のinhouseバージョン作りますが、よろしいですか?",
		Actions: actions,
		Color: "#FFD700",
		CallbackID: "callback_inhouse_deploy",
	}

	_, _, _, _ = api.SendMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
}


func sendMsg(mes string, channel string) {
	_, _, _ = api.PostMessage(channel, slack.MsgOptionText(mes, false))
}