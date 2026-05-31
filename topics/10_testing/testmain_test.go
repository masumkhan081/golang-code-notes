// testmain_test.go
// TestMain — package-level setup and teardown.
//
// If a test file defines TestMain, the testing framework calls it
// instead of running tests directly. You're responsible for invoking
// m.Run() and calling os.Exit with its return code.
//
// Use for: spinning up a test database, starting a mock server,
// seeding fixtures, leak detection (uber-go/goleak.VerifyTestMain).
//
// You may have AT MOST ONE TestMain per package.
package main

import (
	"fmt"
	"os"
	"testing"
)

var testDB string // pretend this is a *sql.DB or similar

func TestMain(m *testing.M) {
	// ---- setup ----
	fmt.Println("setup: starting fake test DB")
	testDB = "connected"

	// ---- run tests ----
	code := m.Run()

	// ---- teardown ----
	// Runs even if tests fail (m.Run returns non-zero, doesn't panic).
	// Does NOT run on panic — wrap in defer if you need that.
	fmt.Println("teardown: closing fake test DB")
	testDB = ""

	os.Exit(code) // must propagate the exit code or the build won't fail on test failures
}

func TestUsesDB(t *testing.T) {
	if testDB != "connected" {
		t.Fatal("TestMain didn't run")
	}
}

/*
COMMON PATTERNS

  func TestMain(m *testing.M) {
      flag.Parse()                 // if your tests read custom flags
      db = openTestDB()
      defer db.Close()             // NOTE: defer doesn't run after os.Exit;
                                   // do teardown before os.Exit instead.
      os.Exit(m.Run())
  }

  func TestMain(m *testing.M) {
      goleak.VerifyTestMain(m)     // fails the suite if any goroutines leak
  }

GOTCHAS
  - defer in TestMain does NOT run if you call os.Exit. Either use a wrapper
    func that returns int, or run teardown explicitly before os.Exit.
  - Only one TestMain per package. For shared setup across packages, use
    init() functions or build tags.
  - If you skip os.Exit(m.Run()), failed tests will not fail the build.

SHORT-MODE SUPPORT
  func TestSlow(t *testing.T) {
      if testing.Short() {
          t.Skip("skipping in -short mode")
      }
      // expensive test
  }
  Run with:  go test -short ./...
*/
