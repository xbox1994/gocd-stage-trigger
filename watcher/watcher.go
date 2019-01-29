package watcher

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"gocd-stage-trigger/gocdcli"
	"strings"
)

func getPipelineNames() ([]string, error) {
	var (
		pipelineOuterJsonResponseBytes string
		e                              error
		pipelineNames                  []string
	)
	if pipelineOuterJsonResponseBytes, e = gocdcli.GetPipelines(); e != nil {
		return nil, e
	}
	var pipelineOuterJsonResponseArray []interface{}
	json.Unmarshal([]byte(pipelineOuterJsonResponseBytes), &pipelineOuterJsonResponseArray)
	for _, outer := range pipelineOuterJsonResponseArray {
		jq := jsonq.NewQuery(outer)
		inner, _ := jq.Array("pipelines")
		for _, pipeline := range inner {
			jq := jsonq.NewQuery(pipeline)
			name, _ := jq.String("name")
			pipelineNames = append(pipelineNames, name)
		}
	}
	return pipelineNames, nil
}

type Stage struct {
	PipelineName string
	StageName    string
}

func getProtectedStages(pipelineNames []string) ([]*Stage, error) {
	var (
		pipelineConfigJsonResponseString string
		e                                error
		stages                           []*Stage
	)
	for _, pipelineName := range pipelineNames {
		if pipelineConfigJsonResponseString, e = gocdcli.GetPipelineConfig(pipelineName); e != nil {
			return nil, e
		}
		pipelineOuterJsonResponseArray := map[string]interface{}{}
		dec := json.NewDecoder(strings.NewReader(pipelineConfigJsonResponseString))
		dec.Decode(&pipelineOuterJsonResponseArray)
		jq := jsonq.NewQuery(pipelineOuterJsonResponseArray)
		array, _ := jq.Array("stages")
		for _, outer := range array {
			jq := jsonq.NewQuery(outer)
			users, _ := jq.Array("approval", "authorization", "users")
			if len(users) > 0 && users[0] == "admin" {
				stageName, _ := jq.String("name")
				stages = append(stages, &Stage{
					PipelineName: pipelineName,
					StageName:    stageName,
				})
			}
		}
	}
	return stages, nil
}

func Run() error {
	var (
		pipelineNames   []string
		protectedStages []*Stage
		e               error
	)
	if pipelineNames, e = getPipelineNames(); e != nil {
		return e
	}
	if protectedStages, e = getProtectedStages(pipelineNames); e != nil {
		return e
	}
	for _, protectedStage := range protectedStages {
		fmt.Println(protectedStage)
	}
	return nil
}
