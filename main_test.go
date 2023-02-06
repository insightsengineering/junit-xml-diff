package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedResult = `
## Unit Test Performance Difference

| Test Suite | $Status$ | Time on ` + "`" + `main` + "`" + ` | $±Time$ | $±Tests$ | $±Skipped$ | $±Failures$ | $±Errors$ |
|:-----|:----:|:----:|:-----:|:-------:|:--------:|:------:|:------:|
| testsuite1 | 💔 | $1.25$ | $+1.30$ | $+2$ | $+2$ | $0$ | $-2$ |
| testsuite3 | 💀 | $3.25$ | $-3.25$ | $-4$ | $0$ | $-3$ | $0$ |
| testsuite5 | 👶 |  | $+7.16$ | $+3$ | $0$ | $0$ | $+1$ |

<details>
  <summary><b>Additional test case details</b></summary>

| Test Suite | $Status$ | Time on ` + "`" + `main` + "`" + ` | $±Time$ | Test Case |
|:-----|:----:|:----:|:----:|:-----|
| testsuite1 | 💔 | $1.21$ | $+5.96$ | testcase1 |
| testsuite1 | 👶 |  | $+0.55$ | testcase3a |
| testsuite1 | 👶 |  | $+0.55$ | testcase3b |
| testsuite2 | 💚 | $5.15$ | $-5.10$ | testcase4 |
| testsuite2 | 💔 | $0.10$ | $+2.00$ | testcase5 |
| testsuite2 | 💀 | $0.08$ | $-0.08$ | testcase6 |
| testsuite3 | 💀 | $0.25$ | $-0.25$ | testcase40 |
| testsuite3 | 💀 | $0.80$ | $-0.80$ | testcase50 |
| testsuite3 | 💀 | $0.58$ | $-0.58$ | testcase60 |
| testsuite5 | 👶 |  | $+1.15$ | testcase400 |
| testsuite5 | 👶 |  | $+5.11$ | testcase500 |
</details>
`

const expectedResultManyTestSuites = `
## Unit Test Performance Difference

<details>
  <summary><b>Test suite performance difference</b></summary>

| Test Suite | $Status$ | Time on ` + "`" + `main` + "`" + ` | $±Time$ | $±Tests$ | $±Skipped$ | $±Failures$ | $±Errors$ |
|:-----|:----:|:----:|:-----:|:-------:|:--------:|:------:|:------:|
| testsuite1 | 💔 | $1.25$ | $+9.00$ | $0$ | $0$ | $+1$ | $0$ |
| testsuite10 | 💔 | $1.25$ | $+18.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite11 | 💔 | $2.25$ | $+18.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite12 | 💔 | $3.25$ | $+29.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite13 | 💔 | $1.25$ | $+20.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite14 | 💔 | $2.25$ | $+20.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite15 | 💔 | $3.25$ | $+20.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite16 | 💔 | $3.25$ | $+30.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite17 | 💔 | $3.25$ | $+21.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite18 | 💚 | $53.25$ | $-28.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite19 | 💚 | $53.25$ | $-27.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite2 | 💔 | $2.25$ | $+10.00$ | $0$ | $0$ | $-1$ | $0$ |
| testsuite20 | 💚 | $43.25$ | $-16.00$ | $0$ | $+1$ | $0$ | $0$ |
| testsuite21 | 💚 | $43.25$ | $-15.00$ | $0$ | $+1$ | $0$ | $0$ |
| testsuite3 | 💔 | $3.25$ | $+10.00$ | $0$ | $0$ | $-2$ | $0$ |
| testsuite4 | 💔 | $1.25$ | $+13.00$ | $0$ | $0$ | $+2$ | $0$ |
| testsuite5 | 💔 | $2.25$ | $+20.00$ | $0$ | $0$ | $-2$ | $0$ |
| testsuite6 | 💔 | $3.25$ | $+12.00$ | $0$ | $0$ | $-1$ | $0$ |
| testsuite7 | 💔 | $1.25$ | $+15.00$ | $0$ | $0$ | $0$ | $-1$ |
| testsuite8 | 💔 | $2.25$ | $+15.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite9 | 💔 | $3.25$ | $+15.00$ | $0$ | $0$ | $0$ | $+1$ |

</details>
`

const expectedResultManyTestSuitesWithChangedTestCases = `
## Unit Test Performance Difference

<details>
  <summary><b>Test suite performance difference</b></summary>

| Test Suite | $Status$ | Time on ` + "`" + `main` + "`" + ` | $±Time$ | $±Tests$ | $±Skipped$ | $±Failures$ | $±Errors$ |
|:-----|:----:|:----:|:-----:|:-------:|:--------:|:------:|:------:|
| testsuite1 | 💔 | $1.25$ | $+9.00$ | $0$ | $0$ | $+1$ | $0$ |
| testsuite10 | 💔 | $1.25$ | $+18.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite11 | 💔 | $2.25$ | $+18.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite12 | 💔 | $3.25$ | $+29.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite13 | 💔 | $1.25$ | $+20.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite14 | 💔 | $2.25$ | $+20.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite15 | 💔 | $3.25$ | $+20.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite16 | 💔 | $3.25$ | $+30.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite17 | 💔 | $3.25$ | $+21.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite18 | 💚 | $53.25$ | $-28.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite19 | 💚 | $53.25$ | $-27.00$ | $0$ | $0$ | $0$ | $0$ |
| testsuite2 | 💔 | $2.25$ | $+10.00$ | $0$ | $0$ | $-1$ | $0$ |
| testsuite20 | 💚 | $43.25$ | $-16.00$ | $0$ | $+1$ | $0$ | $0$ |
| testsuite21 | 💚 | $43.25$ | $-15.00$ | $0$ | $+1$ | $0$ | $0$ |
| testsuite3 | 💔 | $3.25$ | $+10.00$ | $0$ | $0$ | $-2$ | $0$ |
| testsuite4 | 💔 | $1.25$ | $+13.00$ | $0$ | $0$ | $+2$ | $0$ |
| testsuite5 | 💔 | $2.25$ | $+20.00$ | $0$ | $0$ | $-2$ | $0$ |
| testsuite6 | 💔 | $3.25$ | $+12.00$ | $0$ | $0$ | $-1$ | $0$ |
| testsuite7 | 💔 | $1.25$ | $+15.00$ | $0$ | $0$ | $0$ | $-1$ |
| testsuite8 | 💔 | $2.25$ | $+15.00$ | $0$ | $0$ | $0$ | $+1$ |
| testsuite9 | 💔 | $3.25$ | $+15.00$ | $0$ | $0$ | $0$ | $+1$ |

</details>

<details>
  <summary><b>Additional test case details</b></summary>

| Test Suite | $Status$ | Time on ` + "`" + `main` + "`" + ` | $±Time$ | Test Case |
|:-----|:----:|:----:|:----:|:-----|
| testsuite11 | 💚 | $5.15$ | $-3.00$ | testcase11 |
| testsuite14 | 💚 | $5.15$ | $-4.00$ | testcase14 |
| testsuite2 | 💚 | $5.15$ | $-5.00$ | testcase2 |
| testsuite8 | 💔 | $5.15$ | $+3.00$ | testcase8 |
| testsuite9 | 💔 | $0.25$ | $+9.00$ | testcase9 |
</details>
`

func Test_compareXMLReports(t *testing.T) {
	positiveThreshold = 0.2
	negativeThreshold = 0.2
	err := os.MkdirAll("testdata", os.ModePerm)
	checkError(err)
	compareXMLReports("testdata/old.xml", "testdata/new.xml", "testdata/out.md", "main")
	outputBytes, err := os.ReadFile("testdata/out.md")
	checkError(err)
	assert.Equal(t, string(outputBytes), expectedResult)
}

func Test_compareXMLReportsManyTestSuites(t *testing.T) {
	positiveThreshold = 0.2
	negativeThreshold = 0.2
	err := os.MkdirAll("testdata", os.ModePerm)
	checkError(err)
	compareXMLReports("testdata/old_many_test_suites.xml", "testdata/new_many_test_suites.xml", "testdata/out.md", "main")
	outputBytes, err := os.ReadFile("testdata/out.md")
	checkError(err)
	assert.Equal(t, string(outputBytes), expectedResultManyTestSuites)
}

func Test_compareXMLReportsManyTestSuitesWithChangedTestCases(t *testing.T) {
	positiveThreshold = 0.2
	negativeThreshold = 0.2
	err := os.MkdirAll("testdata", os.ModePerm)
	checkError(err)
	compareXMLReports("testdata/old_many_test_suites.xml", "testdata/new_many_test_suites_with_changed_test_cases.xml", "testdata/out.md", "main")
	outputBytes, err := os.ReadFile("testdata/out.md")
	checkError(err)
	assert.Equal(t, string(outputBytes), expectedResultManyTestSuitesWithChangedTestCases)
}
