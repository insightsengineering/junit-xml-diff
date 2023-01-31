package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type TestSuitesXML struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Time      string     `xml:"time,attr"`
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Skipped   int        `xml:"skipped,attr"`
	Failures  int        `xml:"failures,attr"`
	Errors    int        `xml:"errors,attr"`
	TestCases []TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Time      string   `xml:"time,attr"`
	Name      string   `xml:"name,attr"`
	ClassName string   `xml:"classname,attr"`
}

type TestSuiteDiff struct {
	SuiteStatus  string
	TimeDiff     string
	TestsDiff    string
	SkippedDiff  string
	FailuresDiff string
	ErrorsDiff   string
}

type TestCaseTime struct {
	Time         float32
	SuiteName    string
	ClassName    string
	TestCaseName string
}

type TestCaseDiff struct {
	TestCaseStatus string
	TimeDiff       string
	TestCaseName   string
	ClassName      string
	SuiteName      string
}

type TestReport struct {
	SuiteDiff map[string]TestSuiteDiff
	CaseDiff  map[string]TestCaseDiff
}

func getTestSuites(testSuiteXML TestSuitesXML) (map[string]string, map[string]int,
	map[string]int, map[string]int, map[string]int) {
	testSuiteTimes := make(map[string]string)
	testSuiteTests := make(map[string]int)
	testSuiteSkipped := make(map[string]int)
	testSuiteFailures := make(map[string]int)
	testSuiteErrors := make(map[string]int)
	for _, v := range testSuiteXML.TestSuites {
		testSuiteTimes[v.Name] = v.Time
		testSuiteTests[v.Name] = v.Tests
		testSuiteSkipped[v.Name] = v.Skipped
		testSuiteFailures[v.Name] = v.Failures
		testSuiteErrors[v.Name] = v.Errors
	}
	return testSuiteTimes, testSuiteTests, testSuiteSkipped, testSuiteFailures, testSuiteErrors
}

func getTestCases(testSuiteXML TestSuitesXML) map[string]TestCaseTime {
	// Test case is identified by: testsuitename:testcaseclassname:testcasename
	testCaseTimes := make(map[string]TestCaseTime)
	for _, v := range testSuiteXML.TestSuites {
		testSuiteName := v.Name
		for _, testCase := range v.TestCases {
			testCaseTime, err := strconv.ParseFloat(testCase.Time, 32)
			checkError(err)
			testCaseID := testSuiteName + ":" + testCase.ClassName + ":" + testCase.Name
			// It may happen that there are multiple test cases with the same class name and test case name inside a test suite.
			// Times of such cases are added to each other.
			testCaseTimeEntry, ok := testCaseTimes[testCaseID]
			if ok {
				testCaseTimeEntry.Time += float32(testCaseTime)
				testCaseTimes[testCaseID] = testCaseTimeEntry
			} else {
				testCaseTimes[testCaseID] = TestCaseTime{
					float32(testCaseTime),
					testSuiteName,
					testCase.ClassName,
					testCase.Name,
				}

			}
		}
	}
	return testCaseTimes
}

func formatFloat(number float32) string {
	var plusSign string
	if number > 0 {
		plusSign = "+"
	}
	return fmt.Sprintf("%s%.3f", plusSign, number)
}

func formatInt(number int) string {
	var plusSign string
	if number > 0 {
		plusSign = "+"
	}
	return plusSign + strconv.Itoa(number)
}

func compareTestCases(testSuiteOld TestSuitesXML, testSuiteNew TestSuitesXML) map[string]TestCaseDiff {
	testCaseTimesOld := getTestCases(testSuiteOld)
	testCaseTimesNew := getTestCases(testSuiteNew)
	testCaseTimeDiff := make(map[string]TestCaseDiff)
	// Iterate through new test cases times.
	for k := range testCaseTimesNew {
		// Check if the test case existed previously.
		_, ok := testCaseTimesOld[k]
		if ok {
			// Test case exists both in old and new XML.
			testCaseTimeDiff[k] = TestCaseDiff{
				"",
				formatFloat(testCaseTimesNew[k].Time - testCaseTimesOld[k].Time),
				testCaseTimesNew[k].TestCaseName,
				testCaseTimesNew[k].ClassName,
				testCaseTimesNew[k].SuiteName,
			}
		} else {
			// Test case exists only in new XML.
			testCaseTimeDiff[k] = TestCaseDiff{
				"➕",
				formatFloat(testCaseTimesNew[k].Time),
				testCaseTimesNew[k].TestCaseName,
				testCaseTimesNew[k].ClassName,
				testCaseTimesNew[k].SuiteName,
			}
		}
	}
	// Iterate through old test cases times.
	for k := range testCaseTimesOld {
		_, ok := testCaseTimesNew[k]
		if !ok {
			// Test case exists only in old XML.
			testCaseTimeDiff[k] = TestCaseDiff{
				"➖",
				formatFloat(testCaseTimesOld[k].Time),
				testCaseTimesOld[k].TestCaseName,
				testCaseTimesOld[k].ClassName,
				testCaseTimesOld[k].SuiteName,
			}
		}
	}
	return testCaseTimeDiff
}

func compareTestSuites(testSuiteOld TestSuitesXML, testSuiteNew TestSuitesXML) map[string]TestSuiteDiff {

	testSuiteTimesOld, testSuiteTestsOld, testSuiteSkippedOld,
		testSuiteFailuresOld, testSuiteErrorsOld := getTestSuites(testSuiteOld)
	testSuiteTimesNew, testSuiteTestsNew, testSuiteSkippedNew,
		testSuiteFailuresNew, testSuiteErrorsNew := getTestSuites(testSuiteNew)
	testSuiteDiff := make(map[string]TestSuiteDiff)

	// Iterate through test suites in new XML.
	for _, v := range testSuiteNew.TestSuites {
		testSuiteName := v.Name
		newTime := testSuiteTimesNew[testSuiteName]
		newTimeFloat, err := strconv.ParseFloat(newTime, 32)
		checkError(err)
		newTests := testSuiteTestsNew[testSuiteName]
		newSkipped := testSuiteSkippedNew[testSuiteName]
		newFailures := testSuiteFailuresNew[testSuiteName]
		newErrors := testSuiteErrorsNew[testSuiteName]
		// Check if the test suite existed previously.
		// Keys in all maps returned by getTestSuites are the same,
		// so it is okay to iterate through any of these maps.
		_, ok := testSuiteTimesOld[testSuiteName]
		if ok {
			// Test suite name exists both in the old and new XML.
			oldTime := testSuiteTimesOld[testSuiteName]
			oldTimeFloat, err := strconv.ParseFloat(oldTime, 32)
			checkError(err)
			testSuiteDiff[testSuiteName] = TestSuiteDiff{
				"",
				formatFloat(float32(newTimeFloat - oldTimeFloat)),
				formatInt(newTests - testSuiteTestsOld[testSuiteName]),
				formatInt(newSkipped - testSuiteSkippedOld[testSuiteName]),
				formatInt(newFailures - testSuiteFailuresOld[testSuiteName]),
				formatInt(newErrors - testSuiteErrorsOld[testSuiteName]),
			}
		} else {
			// Test suite name exists only in new XML.
			testSuiteDiff[testSuiteName] = TestSuiteDiff{
				"➕",
				formatFloat(float32(newTimeFloat)),
				formatInt(newTests),
				formatInt(newSkipped),
				formatInt(newFailures),
				formatInt(newErrors),
			}
		}
	}
	// Iterate through test suites in old XML.
	for _, v := range testSuiteOld.TestSuites {
		testSuiteName := v.Name
		_, ok := testSuiteTimesNew[v.Name]
		if !ok {
			oldTimeFloat, err := strconv.ParseFloat(testSuiteTimesOld[testSuiteName], 32)
			checkError(err)
			// Test suite name exists only in old XML.
			testSuiteDiff[testSuiteName] = TestSuiteDiff{
				"➖",
				formatFloat(-1 * float32(oldTimeFloat)),
				formatInt(-1 * testSuiteTestsOld[testSuiteName]),
				formatInt(-1 * testSuiteSkippedOld[testSuiteName]),
				formatInt(-1 * testSuiteFailuresOld[testSuiteName]),
				formatInt(-1 * testSuiteErrorsOld[testSuiteName]),
			}
		}
	}
	return testSuiteDiff
}

func compareXMLReports(fileOld, fileNew, fileOut string) {
	xmlFileOld, err := os.Open(fileOld)
	checkError(err)
	xmlFileNew, err := os.Open(fileNew)
	checkError(err)
	defer xmlFileOld.Close()
	defer xmlFileNew.Close()
	byteValueOld, err := io.ReadAll(xmlFileOld)
	checkError(err)
	byteValueNew, err := io.ReadAll(xmlFileNew)
	checkError(err)

	var testSuiteNewXML TestSuitesXML
	var testSuiteOldXML TestSuitesXML
	err = xml.Unmarshal(byteValueNew, &testSuiteNewXML)
	checkError(err)
	err = xml.Unmarshal(byteValueOld, &testSuiteOldXML)
	checkError(err)

	testSuiteDiff := compareTestSuites(testSuiteOldXML, testSuiteNewXML)
	testCasesDiff := compareTestCases(testSuiteOldXML, testSuiteNewXML)

	var testReport TestReport
	testReport.SuiteDiff = testSuiteDiff
	testReport.CaseDiff = testCasesDiff

	tmpl, err := template.New("md").Parse(mdtemplate)
	checkError(err)
	outputFile, err := os.Create(fileOut)
	checkError(err)
	defer outputFile.Close()
	err = tmpl.Execute(outputFile, testReport)
	checkError(err)
}

const mdtemplate = `
| Test Suite | $Status$ | $±Time$ | $±Tests$ | $±Skipped$ | $±Failures$ | $±Errors$ |
|:----:|:----:|:----:|:-----:|:-------:|:--------:|:------:|
{{- range $key, $value := .SuiteDiff }}
| ${{ $key }}$ | ${{ .SuiteStatus }}$ | ${{ .TimeDiff }}$ | ${{ .TestsDiff }}$ | ${{ .SkippedDiff }}$ | ${{ .FailuresDiff }}$ | ${{ .ErrorsDiff }}$ |
{{- end}}

<details>
  <summary><b>Additional test case details</b></summary>

| Test Case | $Status$ | $±Time$ | Test Suite |
|:----:|:----:|:----:|:----:|
{{- range $key, $value := .CaseDiff }}
| {{ .TestCaseName }} | ${{ .TestCaseStatus }}$ | ${{ .TimeDiff }}$ | {{ .SuiteName }} |
{{- end}}
</details>
`

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <old-xml-file-name> <new-xml-file-name> <output-file-name>\n", os.Args[0])
		os.Exit(1)
	}
	compareXMLReports(os.Args[1], os.Args[2], os.Args[3])
}
