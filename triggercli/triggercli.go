package triggercli

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func QueryTriggerStatus(pipelineName string) (string, error) {
	return get("http://10.13.89.38:8080/confirm/query?confirm_type=deploy_k8s&query_key=" + pipelineName)
}

func get(url string) (string, error) {
	client := http.Client{}
	var (
		req *http.Request
		res *http.Response
		e   error
	)
	if req, e = http.NewRequest("GET", url, nil); e != nil {
		return "", e
	}
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
