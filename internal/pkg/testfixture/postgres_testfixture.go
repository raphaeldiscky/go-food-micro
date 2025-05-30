// Package testfixture provides a postgres test fixture.
package testfixture

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	testfixtures "github.com/go-testfixtures/testfixtures/v3"
)

// RunPostgresFixture runs a postgres fixture.
func RunPostgresFixture(db *sql.DB, fixturePaths []string, data map[string]interface{}) error {
	// determine the project's root path
	_, callerPath, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("failed to get caller path")
	}

	// Root folder of this project
	rootPath := filepath.Join(filepath.Dir(callerPath), "../../../..")

	// assemble a list of fixtures paths to be loaded
	for i := range fixturePaths {
		fixturePaths[i] = fmt.Sprintf("%v/%v", rootPath, filepath.ToSlash(fixturePaths[i]))
	}

	// https://github.com/go-testfixtures/testfixtures
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateData(data),
		// Paths must come after Template() and TemplateData()
		testfixtures.Paths(fixturePaths...),
	)
	if err != nil {
		return err
	}

	// load fixtures into DB
	return fixtures.Load()
}
