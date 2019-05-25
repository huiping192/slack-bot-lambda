package bitrise

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)


type BuildStatus int

const (
	Triggered BuildStatus = 0
	Finished BuildStatus = 1
	Aborted BuildStatus = 3
)

type Git struct {
	Provider string `json:"provider"`
	SrcBranch string `json:"src_branch"`
	DstBranch string `json:"dst_branch"`
	PullRequestId int `json:"pull_request_id"`
	Tag string `json:"tag"`
}

type WebHookCallBack struct {
	BuildSlug string `json:"build_slug"`
	BuildNumber int `json:"build_number"`
	AppSlug string `json:"app_slug"`
	BuildStatus int `json:"build_status"`
	BuildTriggeredWorkflow string `json:"build_triggered_workflow"`
	Git Git `json:"message"`
}

type BuildResponse struct {
	Status           int   `json:"status,omitempty"`
	Message string `json:"message"`
	Slug string `json:"slug"`
	Service string `json:"service"`
	BuildSlug string `json:"build_slug"`
	BuildNumber int `json:"build_number"`
	BuildUrl string `json:"build_url"`
	TriggeredWorkflow string `json:"triggered_workflow"`
}

func BuildIap(version string) {
	sendRequest("tag",version,"iap-fastlane")

	branch := "release/" + version
	sendRequest("branch",branch,"iap-fastlane")
}


func BuildInhouse(version string) {
	sendRequest("tag",version,"inhouse-fastlane")

	branch := "release/" + version
	sendRequest("branch",branch,"inhouse-fastlane")
}


func sendRequest(buildType string, buildTarget string, workflow string) BuildResponse {
	token := os.Getenv("BITRISE_TOKEN")
	appId := os.Getenv("BITRISE_APP_ID")

	urlString := "https://app.bitrise.io/app/" + appId +"/build/start.json"

	values := url.Values{}
	values.Add("payload", "{\"hook_info\":{\"type\":\"bitrise\",\"build_trigger_token\":\"" + token + "\"},\"build_params\":{\""+buildType+"\":\""+ buildTarget + "\",\"workflow_id\":\""+ workflow +"\"},\"triggered_by\":\"curl\"}")

	resp, _ := http.PostForm(urlString, values)


	defer resp.Body.Close()

	var responseData BuildResponse
	byteArray, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response" + string(byteArray))

	_ = json.Unmarshal(byteArray, &responseData)

	return responseData
}