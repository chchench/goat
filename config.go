package gotd

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
)

type TestConfig struct {
	Root       string      `json:"testsuite_root"`
	Author     string      `json:"author"`
	Timestamp  string      `json:"modified_timestamp"`
	TestCases  []TestCase  `json:"testcases"`
	TestSuites []TestSuite `json:"testsuites"`
}

type TestCase struct {
	Disabled         bool   `json:"disabled"`
	Name             string `json:"testcase_name"`
	Id               string `json:"testcase_id"`
	Description      string `json:"testcase_description"`
	Kind             string `json:"testcase_kind"`
	AssociatedSuites string `json:"associated_testsuites"`
	Program          string `json:"program"`
	StepsPrerun      string `json:"steps_prerun"`
	StepsPostrun     string `json:"steps_postrun"`
	CliParameters    string `json:"cli_parameters"`
	TimeoutSecs      int    `json:"timeout_secs"`
}

type TestSuite struct {
	Name        string `json:"testsuite_name"`
	Id          string `json:"testsuite_id"`
	Alias       string `json:"testsuite_alias"`
	Description string `json:"testsuite_description"`
}

func (tc *TestConfig) ReadConfigFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error while trying to load test configuration file \"%s\": %v", path, err)
		return err
	}

	err = json.Unmarshal(content, &tc)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return err
	}

	return nil
}

func (tc TestConfig) RetrieveTestCasesByIdentifier(identifier string) []TestCase {

	tcList := []TestCase{}
	for i := 0; i < len(tc.TestCases); i++ {
		match, _ := regexp.MatchString(identifier, tc.TestCases[i].AssociatedSuites)
		if match && !tc.TestCases[i].Disabled {
			tcList = append(tcList, tc.TestCases[i])
		}
	}
	return tcList
}
