package cmd

import (
	"testing"
)

func TestServerCmd(t *testing.T) {
	if serverCmd == nil {
		t.Errorf("serverCmd should not be defined")
	}

	if serverCmd.Use != "server" {
		t.Errorf("serverCmd.Use should be 'server'")
	}

	if serverCmd.Flags().Lookup("port") == nil {
		t.Errorf("expected port flag to be defined")
	}

	if serverCmd.Flags().Lookup("namespace") == nil {
		t.Errorf("expected namespace flag to be defined")
	}
}