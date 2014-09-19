package pipelines

import "time"

type Stage struct {
	PipelineAndStageMixin
}

func (s *Stage) SetDefaults() {
	var err error

	if s.Method == "" {
		s.Method = "POST"
	}

	if s.TimeoutString == "" {
		s.TimeoutString = "1s"
	}

	s.Timeout, err = time.ParseDuration(s.TimeoutString)
	if err != nil {
		s.TimeoutString = "1s"
		s.Timeout, _ = time.ParseDuration(s.TimeoutString)
	}

	s.Accept = "application/json"
	s.ContentType = "application/json"

	s.MergeBodyToChildrenBody()
}
