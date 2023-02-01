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
