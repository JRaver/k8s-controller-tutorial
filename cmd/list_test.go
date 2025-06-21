package cmd

import (
	"testing"
)

func TestListDeploymentCmd(t *testing.T) {
	if listCmd == nil {
		t.Errorf("listCmd should not be defined")
	}

	if listCmd.Use != "list" {
		t.Errorf("listCmd.Use should be 'list'")
	}

	if listCmd.Flags().Lookup("kubeconfig") == nil {
		t.Errorf("expected kubeconfig flag to be defined")
	}

	if listCmd.Flags().Lookup("namespace") == nil {
		t.Errorf("expected namespace flag to be defined")
	}
}