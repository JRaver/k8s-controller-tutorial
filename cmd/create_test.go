package cmd

import (
	"testing"
)

func TestCreateDeploymentCmd(t *testing.T) {
	if createCmd == nil {
		t.Errorf("createCmd should not be defined")
	}

	if createCmd.Use != "create" {
		t.Errorf("createCmd.Use should be 'create'")
	}

	if createCmd.Flags().Lookup("kubeconfig") == nil {
		t.Errorf("expected kubeconfig flag to be defined")
	}

	if createCmd.Flags().Lookup("namespace") == nil {
		t.Errorf("expected namespace flag to be defined")
	}

	if createCmd.Flags().Lookup("deployment-name") == nil {
		t.Errorf("expected deployment-name flag to be defined")
	}
}
