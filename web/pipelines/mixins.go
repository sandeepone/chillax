package pipelines

import (
	"encoding/json"
	"github.com/didip/chillax/libtime"
	"github.com/franela/goreq"
	"io/ioutil"
	"sync"
	"time"
)

type PipelineAndStageMixin struct {
	goreq.Request

	Body   map[string]interface{}
	Stages []*Stage

	// Default is "1s"
	TimeoutString string

	// Default is "5s"
	RetryWaitString string
	FailCount       int

	// Default is 10
	FailMax int
}

type PipelineAndStageSerializableMixin struct {
	Body            map[string]interface{}
	Stages          []*StageSerializable
	TimeoutString   string
	RetryWaitString string
	FailCount       int
	FailMax         int
}

func (mixin *PipelineAndStageMixin) MergeBodyToChildrenBody() {
	for _, stage := range mixin.Stages {
		if stage.Body == nil {
			stage.Body = mixin.Body
		} else {
			for pipelineKey, pipelineValue := range mixin.Body {
				if stage.Body[pipelineKey] == nil {
					stage.Body[pipelineKey] = pipelineValue
				}
			}
		}
	}
}

func (mixin *PipelineAndStageMixin) SetCommonDefaults() {
	if mixin.TimeoutString == "" {
		mixin.TimeoutString = "1s"
	}

	_, err := time.ParseDuration(mixin.TimeoutString)
	if err != nil {
		mixin.TimeoutString = "1s"
	}

	if mixin.RetryWaitString == "" {
		mixin.RetryWaitString = "5s"
	}

	_, err = time.ParseDuration(mixin.RetryWaitString)
	if err != nil {
		mixin.RetryWaitString = "5s"
	}

	if mixin.FailMax <= 0 {
		mixin.FailMax = 10
	}

	if mixin.Method == "" {
		mixin.Method = "POST"
	}

	mixin.Accept = "application/json"
	mixin.ContentType = "application/json"

	mixin.MergeBodyToChildrenBody()
}

func (mixin *PipelineAndStageMixin) Run() RunInstance {
	var err error

	runInstance := mixin.NewRunInstance()

	if mixin.Uri != "" {
		response, err := mixin.Do()

		if err == nil && response != nil && response.Body != nil {
			responseBytes, err := ioutil.ReadAll(response.Body)

			if err == nil && len(responseBytes) > 0 {
				runInstance.ResponseBody = string(responseBytes)
			}
		}

		if err != nil {
			mixin.FailCount += 1

			runInstance.ErrorMessage = err.Error()

			// Retry after sleeping as long as mixin.RetryWaitString
			// only if mixin.FailMax is not exceeded.
			if mixin.FailCount < mixin.FailMax {
				go func(mixin *PipelineAndStageMixin, runInstance RunInstance) {
					libtime.SleepString(mixin.RetryWaitString)

					runInstance = mixin.Run()

				}(mixin, runInstance)
			}

		}
	}

	if err == nil && len(mixin.Stages) > 0 {
		runInstance.RunInstances = make([]RunInstance, len(mixin.Stages))

		var wg sync.WaitGroup

		for i, stage := range mixin.Stages {
			wg.Add(1)

			go func(runInstance RunInstance, i int, stage *Stage) {
				defer wg.Done()

				// Merge the JSON body of previous stage to next stage.
				if len(runInstance.ResponseBody) > 0 {
					json.Unmarshal([]byte(runInstance.ResponseBody), &stage.Body)
				}

				runInstance.RunInstances[i] = stage.Run()
			}(runInstance, i, stage)
		}
		wg.Wait()
	}

	runInstance.Save()

	return runInstance
}

func (mixin *PipelineAndStageMixin) NewRunInstance() RunInstance {
	ri := RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()

	return ri
}
