package gocdcli

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var username = "dev"
var password = "dev"
var host = "http://10.13.89.40:8153"

func Login() error {
	client := http.Client{}
	var (
		req *http.Request
		res *http.Response
		e   error
	)
	if req, e = http.NewRequest("GET", host+"/go/api/current_user", nil); e != nil {
		return e
	}
	req.Header.Add("Accept", "application/vnd.go.cd.v1+json")
	req.SetBasicAuth(username, password)
	if res, e = client.Do(req); e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		//cookie = res.Cookies()[0].String()
		//fmt.Println(cookie)
		return nil
	}
	return errors.New("login failed")
}

func GetPipelines() (string, error) {
	return get(host+"/go/api/config/pipeline_groups", "application/vnd.go.cd.v6+json")
}

func GetPipelineConfig(pipelineName string) (string, error) {
	return get(host+"/go/api/admin/pipelines/"+pipelineName, "application/vnd.go.cd.v6+json")
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
