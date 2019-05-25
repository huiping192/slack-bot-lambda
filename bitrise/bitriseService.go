package bitrise

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Response struct {
	Status           string   `json:"status,omitempty"`
	Message string `json:"message"`
	Slug string `json:"slug"`
	Service string `json:"service"`
	BuildSlug string `json:"build_slug"`
	BuildNumber string `json:"build_number"`
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


func sendRequest(buildType string, buildTarget string, workflow string) Response {
	token := os.Getenv("BITRISE_TOKEN")
	appId := os.Getenv("BITRISE_APP_ID")

	urlString := "https://app.bitrise.io/app/" + appId +"/build/start.json"

	values := url.Values{}
	values.Add("payload", "{\"hook_info\":{\"type\":\"bitrise\",\"build_trigger_token\":\"" + token + "\"},\"build_params\":{\""+buildType+"\":\""+ buildTarget + "\",\"workflow_id\":\""+ workflow +"\"},\"triggered_by\":\"curl\"}")

	resp, _ := http.PostForm(urlString, values)


	defer resp.Body.Close()

	var responseData Response
	byteArray, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response" + string(byteArray))

	_ = json.Unmarshal(byteArray, &responseData)

	return responseData
}