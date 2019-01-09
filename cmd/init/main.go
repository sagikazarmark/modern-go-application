//+build init

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob" // nolint: goimports
	"gopkg.in/AlecAivazis/survey.v1" // nolint: goimports
)

const (
	relativeBase        = "../../"
	originalPackageName = "github.com/sagikazarmark/modern-go-application"
	originalBinaryName  = "modern-go-application"
)

func main() {
	packageName := askPackageName()
	projectName := askProjectName(packageName)

	var questions = []*survey.Question{
		{
			Name: "binaryName",
			Prompt: &survey.Input{
				Message: "Binary name",
				Default: projectName,
				Help:    "Name of the binary package and the built binary",
			},
			Validate: survey.Required,
		},
		{
			Name: "removeInit",
			Prompt: &survey.Confirm{
				Message: "Remove this init script?",
				Default: false,
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		BinaryName string
		RemoveInit bool
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	renameIdeaProjectFile(projectName)
	updateIdeaRunConfigurations(projectName)
	renameBinaryFolder(answers.BinaryName)
	updateMakefile(packageName, answers.BinaryName)
	updateFiles(packageName)
	updateSourceFiles(packageName, answers.BinaryName)
	removeInit(answers.RemoveInit)
}

func fail(err error) {
	if err != nil {
		fmt.Println(err.Error())

		os.Exit(1)
	}
}

func replaceInFile(path string, old string, new string) {
	oldContent, err := ioutil.ReadFile(path)
	fail(err)

	newContent := strings.Replace(string(oldContent), old, new, -1)
	if string(oldContent) == newContent {
		return
	}

	err = ioutil.WriteFile(path, []byte(newContent), 0)
	fail(err)
}

func renameIdeaProjectFile(projectName string) {
	newPathProjectFilePath := fmt.Sprintf(".idea/%s.iml", projectName)
	err := os.Rename(
		filepath.Join(relativeBase, ".idea/project.iml"),
		filepath.Join(relativeBase, newPathProjectFilePath),
	)
	fail(err)

	replaceInFile(filepath.Join(relativeBase, ".idea/modules.xml"), ".idea/project.iml", newPathProjectFilePath)
}

func updateIdeaRunConfigurations(projectName string) {
	runConfigurationFiles := []string{
		"All_tests.xml",
		"Debug.xml",
		"Integration_tests.xml",
		"Tests.xml",
	}

	for _, file := range runConfigurationFiles {
		replaceInFile(
			filepath.Join(relativeBase, ".idea/runConfigurations", file),
			"name=\"project\"",
			fmt.Sprintf("name=\"%s\"", projectName),
		)
	}
}

func renameBinaryFolder(binaryName string) {
	new, old := filepath.Join(relativeBase, "cmd/", originalBinaryName), filepath.Join(relativeBase, "cmd/", binaryName)
	if new == old {
		return
	}

	err := os.Rename(new, old)
	fail(err)
}

func updateMakefile(packageName string, binaryName string) {
	replaceInFile(
		filepath.Join(relativeBase, "Makefile"),
		"PACKAGE = $(shell echo $${PWD\\#\\#*src/})",
		fmt.Sprintf("PACKAGE = %s", packageName),
	)

	replaceInFile(
		filepath.Join(relativeBase, "Makefile"),
		"BUILD_PACKAGE ?= ${PACKAGE}/cmd/$(shell basename $$PWD)",
		fmt.Sprintf("BUILD_PACKAGE = ${PACKAGE}/cmd/%s", binaryName),
	)

	replaceInFile(
		filepath.Join(relativeBase, "Makefile"),
		"BINARY_NAME ?= $(shell basename $$PWD)",
		fmt.Sprintf("BINARY_NAME ?= %s", binaryName),
	)
}

func updateFiles(packageName string) {
	files := []string{
		".circleci/config.yml",
		".gitlab-ci.yml",
		"CHANGELOG.md",
		"Dockerfile",
	}

	for _, file := range files {
		replaceInFile(
			filepath.Join(relativeBase, file),
			originalPackageName,
			packageName,
		)
	}
}

func updateSourceFiles(packageName string, binaryName string) {
	matches, err := zglob.Glob(filepath.Join(relativeBase, "internal/**/*.go"))
	fail(err)

	for _, file := range matches {
		replaceInFile(file, originalPackageName, packageName)
	}

	matches, err = zglob.Glob(filepath.Join(relativeBase, "cmd", binaryName, "/*.go"))
	fail(err)

	for _, file := range matches {
		replaceInFile(file, originalPackageName, packageName)
	}
}

func removeInit(remove bool) {
	err := os.Chdir(relativeBase)
	fail(err)

	if remove {
		err := os.RemoveAll("./cmd/init")
		fail(err)
	}
}
