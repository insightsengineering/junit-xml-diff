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
	XMLName  xml.Name `xml:"testsuite"`
	Time     string   `xml:"time,attr"`
	Name     string   `xml:"name,attr"`
	Tests    int      `xml:"tests,attr"`
	Skipped  int      `xml:"skipped,attr"`
	Failures int      `xml:"failures,attr"`
	Errors   int      `xml:"errors,attr"`
}

type TestSuiteDiff struct {
	SuiteStatus  string
	TimeDiff     string
	TestsDiff    string
	SkippedDiff  string
	FailuresDiff string
	ErrorsDiff   string
}

func getTestSuitesTime(testSuiteXML TestSuitesXML) (map[string]string, map[string]int,
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

func compareTestSuites(testSuiteOld []byte, testSuiteNew []byte) map[string]TestSuiteDiff {
	var testSuiteNewXML TestSuitesXML
	var testSuiteOldXML TestSuitesXML
	err := xml.Unmarshal(testSuiteNew, &testSuiteNewXML)
	checkError(err)
	err = xml.Unmarshal(testSuiteOld, &testSuiteOldXML)
	checkError(err)
	testSuiteTimesOld, testSuiteTestsOld, testSuiteSkippedOld,
		testSuiteFailuresOld, testSuiteErrorsOld := getTestSuitesTime(testSuiteOldXML)
	testSuiteTimesNew, testSuiteTestsNew, testSuiteSkippedNew,
		testSuiteFailuresNew, testSuiteErrorsNew := getTestSuitesTime(testSuiteNewXML)
	testSuiteDiff := make(map[string]TestSuiteDiff)

	// Iterate through test suites in new XML.
	for _, v := range testSuiteNewXML.TestSuites {
		testSuiteName := v.Name
		newTime := testSuiteTimesNew[testSuiteName]
		newTimeFloat, err := strconv.ParseFloat(newTime, 32)
		checkError(err)
		newTests := testSuiteTestsNew[testSuiteName]
		newSkipped := testSuiteSkippedNew[testSuiteName]
		newFailures := testSuiteFailuresNew[testSuiteName]
		newErrors := testSuiteErrorsNew[testSuiteName]
		// Check if the test suite existed previously.
		// Keys in all maps returned by getTestSuitesTimes are the same,
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
	for _, v := range testSuiteOldXML.TestSuites {
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

const mdtemplate = `
|Name|Status|Time|Tests|Skipped|Failures|Errors|
|----|------|----|-----|-------|--------|------|
{{- range $key, $value := .}}
|{{ $key }}|{{ $value.SuiteStatus }}|{{ $value.TimeDiff }}|{{ $value.TestsDiff }}|{{ $value.SkippedDiff }}|{{ $value.FailuresDiff }}|{{ $value.ErrorsDiff }}|
{{- end}}
`

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: ./%s <old-xml-file-name> <new-xml-file-name>\n", os.Args[0])
		os.Exit(1)
	}
	xmlFileOld, err := os.Open(os.Args[1])
	checkError(err)
	xmlFileNew, err := os.Open(os.Args[2])
	checkError(err)
	defer xmlFileOld.Close()
	defer xmlFileNew.Close()
	byteValueOld, err := io.ReadAll(xmlFileOld)
	checkError(err)
	byteValueNew, err := io.ReadAll(xmlFileNew)
	checkError(err)

	testSuiteDiff := compareTestSuites(byteValueOld, byteValueNew)

	tmpl, err := template.New("md").Parse(mdtemplate)
	checkError(err)
	err = tmpl.Execute(os.Stdout, testSuiteDiff)
	checkError(err)
}
