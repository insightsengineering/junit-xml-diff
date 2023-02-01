# junit-xml-diff

[![GitHub Releases](https://img.shields.io/github/v/release/insightsengineering/junit-xml-diff)](https://github.com/insightsengineering/junit-xml-diff/releases)
[![go.mod](https://img.shields.io/github/go-mod/go-version/insightsengineering/junit-xml-diff)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/insightsengineering/junit-xml-diff)](https://goreportcard.com/report/github.com/insightsengineering/junit-xml-diff)
[![Lint Code Base](https://github.com/insightsengineering/junit-xml-diff/actions/workflows/lint.yml/badge.svg)](https://github.com/insightsengineering/junit-xml-diff/actions/workflows/lint.yml)
[![Tests](https://github.com/insightsengineering/junit-xml-diff/actions/workflows/test.yml/badge.svg)](https://github.com/insightsengineering/junit-xml-diff/actions/workflows/test.yml)

Tool for comparing JUnit XML reports

## Usage

### Binary

You may download the binary for your OS/distribution from [https://github.com/insightsengineering/junit-xml-diff/releases](https://github.com/insightsengineering/junit-xml-diff/releases) and run the following command:

```bash
# Assuming you're in the directory which contains the binary
./junit-xml-diff <old.xml> <new.xml> <output.md> <branch-name-corresponding-to-old-xml> <positive-threshold> <negative-threshold>
# For help, run:
./junit-xml-diff
```

### Container

You can pull a container image for this tool from [https://github.com/insightsengineering/junit-xml-diff/pkgs/container/junit-xml-diff](https://github.com/insightsengineering/junit-xml-diff/pkgs/container/junit-xml-diff).

## License

The project is licensed under the [Apache License, Version 2.0](LICENSE).
