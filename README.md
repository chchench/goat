# About GOAT

Go Automated Testing (GOAT) is a testing framework developed in Golang to support automated testing. Its design and development are influenced by the following requirements:

* Enables testing of APIs and SDKs developed for multiple programming languages, including but not limited to JavaScript/TypeScript, Python, Golang, etc.
* Supports configurable tests and test suites so different test cases can be executed for different verification scopes.
* Testing can be invoked and run easily via CLI, and deployed as easily on a developer's laptop or in a centralized, shared testing environment. 
* Supports simple performance benchmarking by recording start time and completion time.
* Supports timing out, to protect against poorly designed or stuck test cases.


## Major Components of GOAT

### The "runtest" Test Driver (TD) Program

This is the command for initiating a new test run to carry out execution of a test suite. This program is essentially the test driver for GOAT. Currently, it supports the following optional command-line flags:

* `-target_testsuite <test suite ID, e.g. "complete">`
	
* `-testconfig string <path to JSON-formatted config file>`
	
	Path to JSON file describing the test cases to execute (default "testcases/testcases.json")

* `-testrun_json_directory <directory path>`
	
	Path of the directory where test run output will be stored. If not specified, working directory will be used.

* `-testrun_json_filename <filename>`
	
	Filename that should be used for outputting test run results. Useful if a fixed filename is needed so, for example, another program can easily locate and process the output.
    	

### Test Configuration File (TCF)

This JSON formatted file describes all the different test cases that can be run. The following is a sample configuration file.

* `testsuite_root` This key has a corresponding value of type string that describes the absolute or relative location of the files/programs associated with each test case. If the value is an absolute path such as `"/a/b/c"` then the complete path for a program identified as `"x/y/sample.exe"` will be `"/a/b/c/x/y/sample.exe"`. If the value is a relative path, then it will be relative to the location of the Test Configuration File. For example, if the value is `"a/b/c"` then the complete path for the aforementioned program will be `"<parent dir of TCF file>/a/b/c/x/y/sample.exe"`.

* `author` This is an optional key with corresponding string value and is currently used for commenting.

* `modified_timestamp` This is an optional key with corresponding string value, and can be used for versioning.

* `testcases` This is an array containing zero or more test cases.

The keys associated with each test case include:

* `testcase_name` This is a string key that can be used to specify a longer, more descriptive identifier for the test case.

* `testcase_id` This is a string key whose value should be unique, and it should be no longer than 8 characters.

* `testcase_description` This string field is for detailed test case description.

* `testcase_kind` This key is not in use right now by the test driver.

* `associated_testsuites` Each test case can be associated with one or more test suites. Use this key-value pair to specify the test suites the test case should be grouped into. Test suite identifiers should be separated using commas.

* `program` This is the path of the main program to be executed as part of test case.

* `cli_parameters` This is an optional key-value for specifying the command-line parameters that should be passed to the program specified by the `program` key-value.

* `steps_prerun` This is an optional key-value pair for specifying the command that need to be executed *prior to* execution of the main test program specified by the `program` key-value. This key-value pair can be used to execute a `make build` command to build the main program or programs required for the test case, for example.

* `steps_postrun` This is an optional key-value pair for specifying the command that need to be executed *after* execution of the main test program specified by the `program` key-value. This key-value pair can be used to execute a `make clean` command to build the main program or programs required for the test case, for example.

* `timeout_secs` Maximum amount of time allowed for execution of the test case in seconds. The value specified should be an integer.

* `disabled` If this key-value pair exists and has a value of `true`, then the test case will not be included in any test suite (and test run).

At this time, the `testsuites` portion of the configuration file is not used by the Test Driver and can be ignored.





	{
		"testsuite_root": "",
		"author": "Charles",
		"modified_timestamp": "",
		"testcases": [
			{
				"testcase_name": "tc-001-detect",
				"testcase_id": "tc-001",
				"testcase_description": "Checks to make sure target files can be detected",
	            "testcase_kind": "__NOT_USED_CURRENTLY__",
				"associated_testsuites": "smoke,complete",
	            "program": "golang/tc-001-detect/scan.exe",
	            "cli_parameters": "-dir /tmp/testsamples/foo",
				"steps_prerun": "make -C golang/tc-001-detect build",
	            "steps_postrun": "make -C golang/tc-001-detect clean",
	            "timeout_secs": 300,
	            "disabled": true
			},
			{
				... similar entry to the "tc-001-detect" example above ...
			}
		]
	}
	


### Sample Test Directory Hierarchy

```
test-root
│   test-configuration-file.json    
│
└───golang
│   │   
│   └───testcase-1
│       │   file111.txt
│       │   file112.txt
│       │   ...
│   
└───typescript
    │   file021.txt
    │   file022.txt
```

The Test Configuration File (TCF) can reside out of the directory tree above, but in that situation the value for the top-level key of `testsuite_root` inside TCF should probably be specified so it is clear to everyone where all the different test cases are stored.


## Getting Started

### Prerequisites

The following instruction assumes your target environment is an Ubuntu Linux environment. It hasn't been verified for other operating systems.


### Installation

__TBD__

## Usage

__TBD__


## Support

__TBD__
