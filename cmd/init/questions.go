package main

import (
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/AlecAivazis/survey.v1"
)

func mustAskOne(p survey.Prompt, response interface{}, v survey.Validator, opts ...survey.AskOpt) {
	err := survey.AskOne(p, response, v, opts...)
	fail(err)
}

// nolint: gochecknoglobals
var defaultPackageRegex = regexp.MustCompile("^(?:/.+)+/src/(.*)/cmd/init$")

func askPackageName() string {
	var defaultPackageName string

	cwd, err := os.Getwd()
	if err == nil {
		if matches := defaultPackageRegex.FindStringSubmatch(cwd); len(matches) > 1 {
			defaultPackageName = matches[1]
		}
	}

	var packageName string

	const (
		name = "Package name"
		help = "Go package/module name"
	)

	mustAskOne(&survey.Input{Message: name, Default: defaultPackageName, Help: help}, &packageName, survey.Required)

	return packageName
}

func askProjectName(packageName string) string {
	defaultProjectName := filepath.Base(packageName)

	var projectName string

	const (
		name = "Project name"
		help = "Short project name used in IDE config"
	)

	mustAskOne(&survey.Input{Message: name, Default: defaultProjectName, Help: help}, &projectName, survey.Required)

	return projectName
}
