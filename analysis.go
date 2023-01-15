package gotd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/TwiN/go-color"
)

func CompareRuns(path_target_run string, path_comp_run string) error {
	tr1 := readTestRunFile(path_comp_run)
	tr2 := readTestRunFile(path_target_run)

	log.Printf("TARGET RUN:  %s (%v test cases executed over %s)\n", tr2.RunID, len(tr2.ExecutedTestResults), tr2.GetExecutionDuration())
	log.Printf("COMP RUN:    %s (%v test cases executed over %s)\n", tr1.RunID, len(tr1.ExecutedTestResults), tr1.GetExecutionDuration())
	log.Printf("\n")

	var testCaseResultMap = make(map[string]*TestCaseResult)

	for i := 0; i < len(tr1.ExecutedTestResults); i++ {
		testCaseResultMap[tr1.ExecutedTestResults[i].Id] = &tr1.ExecutedTestResults[i]
	}

	formatting := "%15s   %10s   %12s   %15s   %10s   %12s   %12s\n"
	log.Printf(formatting+"\n", "TARGET RUN", "P/F/T", "TIME", "COMP RUN", "P/F/T", "TIME", "TIME DELTA")
	log.Printf("")

	var totalDuration4PriorRun time.Duration

	var numSharedCases int
	var totalSharedCaseTime4Target time.Duration
	var totalSharedCaseTime4Comp time.Duration

	for i := 0; i < len(tr2.ExecutedTestResults); i++ {

		tcr := tr2.ExecutedTestResults[i]
		tcr_comp := testCaseResultMap[tcr.Id]

		if tcr_comp != nil {

			totalDuration4PriorRun += tcr_comp.GetExecutionDuration()

			timeDiff := tcr.GetExecutionDuration() - tcr_comp.GetExecutionDuration()

			log.Printf(formatting, tcr.Id, getTestCaseResult(tcr), tcr.GetExecutionDuration(),
				tcr_comp.Id, getTestCaseResult(*tcr_comp), tcr_comp.GetExecutionDuration(), timeDiff)

			numSharedCases++
			totalSharedCaseTime4Target += tcr.GetExecutionDuration()
			totalSharedCaseTime4Comp += tcr_comp.GetExecutionDuration()
		} else {
			log.Printf(formatting, tcr.Id, getTestCaseResult(tcr), tcr.GetExecutionDuration(),
				"N/A", "N/A", "N/A", "N/A")
		}
	}

	sharedCaseTimeDiff := totalSharedCaseTime4Target - totalSharedCaseTime4Comp

	log.Printf("\n")
	log.Printf("For the %v common test cases executed, TARGET took %s and COMP took %s; ", numSharedCases, totalSharedCaseTime4Target, totalSharedCaseTime4Comp)
	log.Printf("delta = %v (%6.2f%%)\n", sharedCaseTimeDiff, float64(sharedCaseTimeDiff*100)/float64(totalSharedCaseTime4Comp))

	return nil
}

func readTestRunFile(path string) TestRun {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var tr TestRun
	err = json.Unmarshal(content, &tr)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return tr
}

func getTestCaseResult(tcr TestCaseResult) string {
	if tcr.TimedOut {
		return "TTTTT"
	}
	if tcr.Passed {
		return "+++++"
	}
	return "-----"
}

func getTestCaseResultColor(tcr TestCaseResult) string {
	if tcr.TimedOut {
		return (color.Colorize(color.Yellow, "TTTTT"))
	}
	if tcr.Passed {
		return (color.Colorize(color.Green, "+++++"))
	}
	return (color.Colorize(color.Red, "-----"))
}
