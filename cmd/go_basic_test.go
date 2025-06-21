package cmd

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	k := &Kubernetes{
		name:       "test",
		version:    "1.0.0",
		users:      []string{"user1", "user2"},
		nodeNumber: 1,
	}

	// Check initial count of users
	initialCount := len(k.users)
	if initialCount != 2 {
		t.Errorf("Expected 2 users, got %d", initialCount)
	}

	k.AddUser("user3")

	// Check count of users after adding
	finalCount := len(k.users)
	if finalCount != 3 {
		t.Errorf("Expected 3 users after adding, got %d", finalCount)
	}

	// Check if the correct user was added
	if k.users[2] != "user3" {
		t.Errorf("Expected 'user3', got '%s'", k.users[2])
	}
}

func TestRemoveUser(t *testing.T) {
	k := &Kubernetes{
		name:       "test",
		version:    "1.0.0",
		users:      []string{"user1", "user2", "user3"},
		nodeNumber: 1,
	}

initialCount := len(k.users)
	if initialCount != 3 {
		t.Errorf("Expected 3 users, got %d", initialCount)
	}

	k.RemoveUser("user2")

	// Check count of users after removing
	finalCount := len(k.users)
	if finalCount != 2 {
		t.Errorf("Expected 2 users after removing, got %d", finalCount)
	}
}

func TestGoBasicCmd(t *testing.T) {
	if goBasicCmd == nil {
		t.Errorf("serverCmd should be defined")
	}

	if goBasicCmd.Use != "go-basic" {
		t.Errorf("serverCmd.Use should be 'go-basic'")
	}
	if goBasicCmd.PersistentFlags() == nil {
		t.Errorf("Command go-basic should haven't any flags")
	}

}