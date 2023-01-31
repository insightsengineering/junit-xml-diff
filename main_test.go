package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedResult = `
|Test Suite|Status|±Time|±Tests|±Skipped|±Failures|±Errors|
|----|----|----|-----|-------|--------|------|
|testsuite1||+0.300|+2|+2|0|-2|
|testsuite2||0.000|-1|+2|-2|0|
|testsuite3|➖|-3.250|-4|0|-3|0|
|testsuite5|➕|+7.158|+3|0|0|+1|

<details>
  <summary><b>Additional test case details</b></summary>

|Test Case|Status|±Time|Test Suite|
|----|----|----|----|
|testcase1||+5.960|testsuite1|
|testcase2||0.000|testsuite1|
|testcase3||0.000|testsuite1|
|testcase3a|➕|+0.550|testsuite1|
|testcase3b|➕|+0.550|testsuite1|
|testcase4||-0.100|testsuite2|
|testcase5||+2.000|testsuite2|
|testcase6|➖|+0.080|testsuite2|
|testcase40|➖|+0.250|testsuite3|
|testcase50|➖|+0.800|testsuite3|
|testcase60|➖|+0.580|testsuite3|
|testcase400|➕|+1.150|testsuite5|
|testcase500|➕|+5.111|testsuite5|
</details>
`

func Test_compareXMLReports(t *testing.T) {
	err := os.MkdirAll("testdata", os.ModePerm)
	checkError(err)
	compareXMLReports("testdata/old.xml", "testdata/new.xml", "testdata/out.md")
	outputBytes, err := os.ReadFile("testdata/out.md")
	assert.Equal(t, string(outputBytes), expectedResult)
}
