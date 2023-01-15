package gotd

import (
	"encoding/json"
	"time"
)

type TestCaseResult struct {
	Id        string    `json:"testcase_result_id"`
	StartTime time.Time `json:"testcase_result_starttime"`
	EndTime   time.Time `json:"testcase_result_endtime"`
	CmdLine   string    `json:"testcase_result_cmdline"`
	Passed    bool      `json:"testcase_result_passed"`
	TimedOut  bool      `json:"testcase_result_timed_out"`
	Data      string    `json:"testcase_result_data"`
	ParentId  string    `json:"testcase_parent_id"`
}

func (tc *TestCaseResult) StartTest() {
	tc.StartTime = time.Now()
}

func (tc *TestCaseResult) EndTest() {
	tc.EndTime = time.Now()
}

func (tc TestCaseResult) GetExecutionDuration() time.Duration {
	d := tc.EndTime.Sub(tc.StartTime)
	return d
}

func (tc TestCaseResult) GetPassFailStatus() bool {
	return tc.Passed
}

func (tc *TestCaseResult) Marshal2JSON() (string, error) {
	result_str, err := json.Marshal(*tc)
	return string(result_str), err
}

func (tc *TestCaseResult) UnmarshalFromJSON(result_str string) error {
	b := []byte(result_str)
	return json.Unmarshal(b, tc)
}
