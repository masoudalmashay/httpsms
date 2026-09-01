package di

import (
	"os"
	"os/exec"
	"testing"
)

func TestAxiomLoggerFallsBackToConsoleWhenConfigIsMissing(t *testing.T) {
	if os.Getenv("TEST_AXIOM_LOGGER_CHILD") == "1" {
		t.Setenv("ENV", "production")
		t.Setenv("AXIOM_TOKEN", "")
		t.Setenv("AXIOM_DATASET_EVENTS", "")

		logger := axiomLogger(3)
		if logger == nil {
			t.Fatal("expected a logger instance")
		}
		if logger.Logger == nil {
			t.Fatal("expected an underlying zerolog logger")
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestAxiomLoggerFallsBackToConsoleWhenConfigIsMissing")
	cmd.Env = append(os.Environ(), "TEST_AXIOM_LOGGER_CHILD=1")
	if err := cmd.Run(); err != nil {
		t.Fatalf("axiomLogger should fall back to console logging when config is missing: %v", err)
	}
}
