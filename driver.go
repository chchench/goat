package gotd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TwiN/go-color"
)

type TestRun struct {
	RunID               string           `json:"testrun_id"`
	Config              string           `json:"testrun_config"`
	StartTime           time.Time        `json:"testrun_startime"`
	EndTime             time.Time        `json:"testrun_endtime"`
	TestCases           []TestCase       `json:"testrun_testcases"`
	ExecutedTestResults []TestCaseResult `json:"testrun_executed_test_results"`
}

func (tr *TestRun) GetExecutionDuration() time.Duration {
	d := tr.EndTime.Sub(tr.StartTime)
	return d
}

type TestRunAnalysis struct {
	NumTCExecuted    int
	NumSubTCExecuted int
	NumPassed        int
	NumFailed        int
	NumTimedOut      int
	TotalDuration    time.Duration
}

func (tr *TestRun) GenerateTestRunAnalysis() TestRunAnalysis {
	var tra TestRunAnalysis

	num := len(tr.ExecutedTestResults)

	for i := 0; i < num; i++ {

		isSubTC := (tr.ExecutedTestResults[i].ParentId != "")
		if isSubTC {
			tra.NumSubTCExecuted++
			continue
		}

		tra.NumTCExecuted++

		if tr.ExecutedTestResults[i].TimedOut {
			tra.NumTimedOut++
		} else if tr.ExecutedTestResults[i].Passed {
			tra.NumPassed++
		} else {
			tra.NumFailed++
		}

		tra.TotalDuration += tr.ExecutedTestResults[i].GetExecutionDuration()
	}

	return tra
}

func (tr *TestRun) DumpTestRun2JSON(filename string, dirpath string) string {

	if dirpath == "" {
		dirpath = "."
	}

	if filename == "" {
		filename = tr.RunID
	}

	path := dirpath + "/" + filename

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
	}

	r, err := json.Marshal(*tr)
	if err == nil {
		if _, err = f.Write(r); err != nil {
			panic(err)
		}
	}

	return string(r)
}

type TestDriver struct {
	TestConfigPath    string
	TestSuiteRootPath string
	config            TestConfig
	testRuns          []*TestRun
	mu                sync.Mutex
}

func (td *TestDriver) LoadTestConfiguration(path string) error {
	err := td.config.ReadConfigFile(path)
	if err == nil {
		abs, _ := filepath.Abs(path)
		td.TestConfigPath = abs
	}

	configFileParentDir := filepath.Dir(td.TestConfigPath)
	if filepath.IsAbs(td.config.Root) {
		td.TestSuiteRootPath, _ = filepath.Abs(td.config.Root)
	} else {
		td.TestSuiteRootPath = filepath.Join(configFileParentDir, td.config.Root)
	}

	return err
}

func (td *TestDriver) createNewTestRun() *TestRun {
	tr := TestRun{}
	td.mu.Lock()
	td.testRuns = append(td.testRuns, &tr)
	td.mu.Unlock()
	return &tr
}

func (td *TestDriver) ExecuteTestSuite(identifier string) TestRun {

	tr := td.createNewTestRun()
	tr.Config = td.TestConfigPath
	tr.TestCases = td.config.RetrieveTestCasesByIdentifier(identifier)

	tr.StartTime = time.Now()

	log.Printf("There are %v test case(s) to be executed based on selection of test suite \"%s\" ...\n",
		len(tr.TestCases), identifier)

	for i := 0; i < len(tr.TestCases); i++ {
		log.Printf("[%04d] Executing test case [%s/%s] ...\n",
			i+1, tr.TestCases[i].Id, tr.TestCases[i].Name)

		if td.config.Root != "" {
			tr.TestCases[i].Program = td.config.Root + "/" + tr.TestCases[i].Program
		}
		tcr, err, subcaseResults := executeTestCase(tr.TestCases[i], td.TestSuiteRootPath)

		var exec_str string
		if err == nil {
			exec_str = fmt.Sprintf("[%04d] %s %s (time of %v for %s)",
				i+1, tcr.Id, getPassedFailedMsg(tcr.Passed), tcr.GetExecutionDuration(), tcr.CmdLine)

		} else if tcr.TimedOut || err.Error() == "timed out" {
			exec_str = fmt.Sprintf("[%04d] %s "+color.Colorize(color.Red,
				"TIMED OUT")+" (terminated after %v for %s)",
				i+1, tcr.Id, tcr.GetExecutionDuration(), tcr.CmdLine)

		} else {
			exec_str = fmt.Sprintf("[%04d] %s "+color.Colorize(color.Red,
				"FAILED TO RUN")+" (time of %v for %s, error = '%v')",
				i+1, tcr.Id, tcr.GetExecutionDuration(), tcr.CmdLine, err)
		}

		log.Println(exec_str)
		tcr.Data = exec_str

		tr.ExecutedTestResults = append(tr.ExecutedTestResults, *tcr)
		if err == nil {
			processSubcaseResults(tcr, tr, &subcaseResults)

			for j := 0; j < len(subcaseResults); j++ {
				scr := subcaseResults[j]
				exec_str = fmt.Sprintf("[%04d-%04d] %s %s (time of %v for %s)",
					i+1, j+1, scr.Id, getPassedFailedMsg(scr.Passed),
					scr.GetExecutionDuration(), scr.CmdLine)
				log.Println(exec_str)
			}
		}
	}

	tr.EndTime = time.Now()

	tr.RunID = fmt.Sprintf("Run-%s-%s-%s",
		identifier, tr.StartTime.Format(time.RFC3339), tr.EndTime.Format(time.RFC3339))

	return *tr
}

func getPassedFailedMsg(passed bool) string {
	if passed {
		return color.Colorize(color.Green, "PASSED")
	}
	return color.Colorize(color.Red, "FAILED")
}

func processSubcaseResults(tcr *TestCaseResult, tr *TestRun, subcaseResults *[]TestCaseResult) {
	for i := 0; i < len(*subcaseResults); i++ {
		(*subcaseResults)[i].Id = fmt.Sprintf("%s:%s", tcr.Id, (*subcaseResults)[i].Id)
		(*subcaseResults)[i].ParentId = tcr.Id
		tr.ExecutedTestResults = append(tr.ExecutedTestResults, (*subcaseResults)[i])
	}
}

type executionResult struct {
	tcr            *TestCaseResult
	err            error
	subcaseResults []TestCaseResult
}

func executeTestCase(tc TestCase, tsRootPath string) (*TestCaseResult, error, []TestCaseResult) {
	tcr := TestCaseResult{}
	tcr.Id = tc.Id

	if tc.StepsPrerun != "" {
		prog, params := parseCmdLine4ProgNParams(tc.StepsPrerun)

		if _, err := LaunchCommand(tsRootPath, prog, params); err != nil {
			return &tcr, fmt.Errorf("test case %s pre-run steps [%s] failed: %s",
				tcr.Id, tc.StepsPrerun, err), nil
		} else {
			log.Printf("Successfully executed pre-run command '%s'\n", tc.StepsPrerun)
		}
	}

	tcr.CmdLine = tc.Program + " " + tc.CliParameters
	result := make(chan executionResult, 1)

	tcr.StartTest()
	go func() {
		result <- executeTestCaseInternal(tsRootPath, tc, &tcr)
	}()

	var ptcr *TestCaseResult
	var err error
	var subcaseResults []TestCaseResult

	select {

	case <-time.After(time.Duration(tc.TimeoutSecs) * time.Second):
		tcr.EndTest()
		tcr.TimedOut = true
		ptcr = &tcr
		err = fmt.Errorf("test case %s timed out", tcr.Id)
		subcaseResults = nil

	case result := <-result:
		ptcr = result.tcr
		err = result.err
		subcaseResults = result.subcaseResults
	}

	if tc.StepsPostrun != "" {
		prog, params := parseCmdLine4ProgNParams(tc.StepsPostrun)
		if _, err := LaunchCommand(tsRootPath, prog, params); err != nil {
			return &tcr, fmt.Errorf("test case %s post-run steps [%s] failed: %s",
				tcr.Id, tc.StepsPostrun, err), nil
		} else {
			log.Printf("Successfully executed post-run command '%s'\n", tc.StepsPostrun)
		}
	}

	return ptcr, err, subcaseResults
}

func executeTestCaseInternal(tsRootPath string, tc TestCase, tcr *TestCaseResult) executionResult {

	tcr.Id = tc.Id

	// TBD: Need to clean up how testcase program command line is passed down for invocation
	// and constructed once there is a more clear understanding of how the test programs will
	// be launched.

	paramList := splitNClean(tc.CliParameters)

	tcr.StartTest()

	log.Printf("Executing command [%s %s] ...", tc.Program, strings.Join(paramList, " "))
	output, err := LaunchCommand(tsRootPath, tc.Program, paramList)

	tcr.EndTest()

	if err != nil {
		return executionResult{
			tcr:            tcr,
			err:            err,
			subcaseResults: nil,
		}
	}

	// TBD: Process stdout to determine result of test case

	err, passed, subcaseResults := parseTestResult(output)
	tcr.Passed = passed

	return executionResult{
		tcr:            tcr,
		err:            err,
		subcaseResults: subcaseResults,
	}
}

type ScanResult struct {
	Id         string    `json:"id"`
	StartTime  time.Time `json:"scan_starttime"`
	EndTime    time.Time `json:"scan_endtime"`
	Test       string    `json:"test"`
	TestPassed bool      `json:"test_passed"`
	Data       string    `json:"data`
}

type OverallTestResult struct {
	StartTime    time.Time    `json:"starttime"`
	EndTime      time.Time    `json:"endtime"`
	Passed       bool         `json:"overall_passed"`
	CombinedData []ScanResult `json:"combined_scan_results"`
}

func parseTestResult(result string) (error, bool, []TestCaseResult) {

	//var m map[string]interface{}

	var resultObj OverallTestResult
	err := json.Unmarshal([]byte(result), &resultObj)

	if err != nil {
		return err, false, nil
	}

	var subCaseResults []TestCaseResult

	for i := 0; i < len(resultObj.CombinedData); i++ {
		scanResult := resultObj.CombinedData[i]

		var tcr TestCaseResult
		tcr.StartTime = scanResult.StartTime
		tcr.EndTime = scanResult.EndTime
		tcr.Passed = scanResult.TestPassed
		tcr.CmdLine = scanResult.Test
		tcr.Data = scanResult.Data
		tcr.Id = scanResult.Id

		subCaseResults = append(subCaseResults, tcr)
	}

	return nil, resultObj.Passed, subCaseResults
}
