package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedResult = `
| Test Suite | $Status$ | Time for ` + "`" + `main` + "`" + ` | $±Time$ | $±Tests$ | $±Skipped$ | $±Failures$ | $±Errors$ |
|:-----|:----:|:----:|:-----:|:-------:|:--------:|:------:|:------:|
| testsuite1 | 💔 | $1.250$ | $+1.300$ | $+2$ | $+2$ | $0$ | $-2$ |
| testsuite2 |  | $2.250$ |  | $-1$ | $+2$ | $-2$ | $0$ |
| testsuite3 | 💀 | $3.250$ | $-3.250$ | $-4$ | $0$ | $-3$ | $0$ |
| testsuite5 | 👶 |  | $+7.158$ | $+3$ | $0$ | $0$ | $+1$ |

<details>
  <summary><b>Additional test case details</b></summary>

| Test Suite | $Status$ | Time for ` + "`" + `main` + "`" + ` | $±Time$ | Test Case |
|:-----|:----:|:----:|:----:|:-----|
| testsuite1 | 💔 | $1.210$ | $+5.960$ | testcase1 |
| testsuite1 |  | $0.320$ |  | testcase2 |
| testsuite1 |  | $0.550$ |  | testcase3 |
| testsuite1 | 👶 |  | $+0.550$ | testcase3a |
| testsuite1 | 👶 |  | $+0.550$ | testcase3b |
| testsuite2 | 💚 | $5.150$ | $-5.100$ | testcase4 |
| testsuite2 | 💔 | $0.100$ | $+2.000$ | testcase5 |
| testsuite2 | 💀 | $0.080$ | $+0.080$ | testcase6 |
| testsuite3 | 💀 | $0.250$ | $+0.250$ | testcase40 |
| testsuite3 | 💀 | $0.800$ | $+0.800$ | testcase50 |
| testsuite3 | 💀 | $0.580$ | $+0.580$ | testcase60 |
| testsuite5 | 👶 |  | $+1.150$ | testcase400 |
| testsuite5 | 👶 |  | $+5.111$ | testcase500 |
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
