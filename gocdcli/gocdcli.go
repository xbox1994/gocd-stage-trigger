package gocdcli

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var username = "dev"
var password = "dev"
var host = "http://10.13.89.40:8153"

func GetPipelines() (string, error) {
	return get(host+"/go/api/config/pipeline_groups", "application/vnd.go.cd.v6+json")
}

func GetPipelineConfig(pipelineName string) (string, error) {
	return get(host+"/go/api/admin/pipelines/"+pipelineName, "application/vnd.go.cd.v6+json")
}

func TriggerStage(pipelineName, stageName, counter string) (string, error) {
	return post(host+"/go/api/stages/"+pipelineName+"/"+counter+"/"+stageName+"/run", "application/vnd.go.cd.v1+json")
}

func get(url, accept string) (string, error) {
	client := http.Client{}
	var (
		req *http.Request
		res *http.Response
		e   error
	)
	if req, e = http.NewRequest("GET", url, nil); e != nil {
		return "", e
	}
	req.Header.Add("Accept", accept)
	req.SetBasicAuth(username, password)
	if res, e = client.Do(req); e != nil {
		return "", e
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return string(bodyBytes), nil
	}
	return "", errors.New(url + "get failed")
}

func post(url, accept string) (string, error) {
	client := http.Client{}
	var (
		req *http.Request
		res *http.Response
		e   error
	)
	if req, e = http.NewRequest("POST", url, nil); e != nil {
		return "", e
	}
	req.Header.Add("Accept", accept)
	req.Header.Add("X-GoCD-Confirm", "true")
	req.SetBasicAuth(username, password)
	if res, e = client.Do(req); e != nil {
		return "", e
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return string(bodyBytes), nil
	}
	return "", errors.New(url + "get failed")
}
