package main

import (
	"flag"
	"log"

	"github.com/chchench/gotd"
)

func main() {

	var testConfigPath string
	var testSuite2Exec string
	var testRunFilename string
	var testRunDirPath string

	flag.StringVar(&testConfigPath, "testconfig", "testcases.json",
		"Path to JSON file describing the test cases to execute")
	flag.StringVar(&testSuite2Exec, "target_testsuite", "complete",
		"Test suite containing the test cases to execute")
	flag.StringVar(&testRunFilename, "testrun_json_filename", "",
		"Name of the test run output file; if not specified, default will be used")
	flag.StringVar(&testRunDirPath, "testrun_json_directory", "",
		"Path of the directory where test run output will be stored; if not specified, default will be working dir")

	flag.Parse()

	testdriver := gotd.TestDriver{}

	err := testdriver.LoadTestConfiguration(testConfigPath)
	if err != nil {
		log.Fatal("Unable to load test configuration file " + testConfigPath)
	}

	log.Printf("Data from %s are loaded\n", testdriver.TestConfigPath)
	log.Printf("Test suite root = %s\n", testdriver.TestSuiteRootPath)

	tr := testdriver.ExecuteTestSuite(testSuite2Exec)
	tr.DumpTestRun2JSON(testRunFilename, testRunDirPath)

	tra := tr.GenerateTestRunAnalysis()

	log.Printf("Total ran: %v / Total Passed: %v / Total Failed: %v / Total Timed-out: %v\n",
		tra.NumTCExecuted, tra.NumPassed, tra.NumFailed, tra.NumTimedOut)
	log.Printf("Total execution time for this test run = %v\n", tra.TotalDuration)
}
