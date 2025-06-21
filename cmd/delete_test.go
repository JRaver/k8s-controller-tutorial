package cmd

import (
	"testing"
)

func TestDeleteDeploymentCmd(t *testing.T) {
	if deleteCmd == nil {
		t.Errorf("deleteCmd should not be defined")
	}

	if deleteCmd.Use != "delete" {
		t.Errorf("deleteCmd.Use should be 'delete'")
	}

	if deleteCmd.Flags().Lookup("kubeconfig") == nil {
		t.Errorf("expected kubeconfig flag to be defined")
	}

	if deleteCmd.Flags().Lookup("namespace") == nil {
		t.Errorf("expected namespace flag to be defined")
	}

	if deleteCmd.Flags().Lookup("deployment-name") == nil {
		t.Errorf("expected deployment-name flag to be defined")
	}
}