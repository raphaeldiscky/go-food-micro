// Package environment provides a module for the environment.
package environment

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"emperror.dev/errors"
	"github.com/spf13/viper"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
)

// FixProjectRootWorkingDirectoryPath fixes the project root working directory path.
func FixProjectRootWorkingDirectoryPath() {
	currentWD, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current working directory is: `%s`", currentWD)

	rootDir := GetProjectRootWorkingDirectory()
	// change working directory
	err = os.Chdir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	newWD, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("New fixed working directory is: `%s`", newWD)
}

// GetProjectRootWorkingDirectory gets the project root working directory.
func GetProjectRootWorkingDirectory() string {
	var rootWorkingDirectory string
	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	// viper will get it from `os env` or a .env file
	pn := viper.GetString(constants.PROJECT_NAME_ENV)
	if pn != "" {
		rootWorkingDirectory = getProjectRootDirectoryFromProjectName(pn)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dir, err := searchRootDirectory(wd)
		if err != nil {
			log.Fatal(err)
		}
		rootWorkingDirectory = dir
	}

	absoluteRootWorkingDirectory, err := filepath.Abs(rootWorkingDirectory)
	if err != nil {
		log.Fatal(err)
	}

	return absoluteRootWorkingDirectory
}

// getProjectRootDirectoryFromProjectName gets the project root directory from the project name.
func getProjectRootDirectoryFromProjectName(pn string) string {
	// set root working directory of our app in the viper
	// https://stackoverflow.com/a/47785436/581476
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for !strings.HasSuffix(wd, pn) {
		wd = filepath.Dir(wd)
	}

	return wd
}

// searchRootDirectory searches the root directory.
func searchRootDirectory(
	dir string,
) (string, error) {
	// List files and directories in the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.WrapIf(err, "Error reading directory")
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.EqualFold(
				fileName,
				"go.mod",
			) {
				return dir, nil
			}
		}
	}

	// If no config file found in this directory, recursively search its parent
	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		// We've reached the root directory, and no go.mod file was found
		return "", errors.WrapIf(err, "No go.mod file found")
	}

	return searchRootDirectory(parentDir)
}
