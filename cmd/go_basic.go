package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var goBasicCmd = &cobra.Command{
	Use:   "go-basic",
	Short: "Basic Go commands for Lab 1",
	Run: func(cmd *cobra.Command, args []string) {
		k := &Kubernetes{
			name: "Kubernetes",
			version: "1.24.0",
			users: []string{"user1", "user2", "user3"},
			nodeNumber: 3,
		}
		// Print users
		fmt.Println(k.PrintUsers())
		// Add user
		k.AddUser("user4")
		fmt.Println(k.PrintUsers())
		// Remove user
		k.RemoveUser("user2")
		fmt.Println(k.PrintUsers())
	},
}

func init() {
	rootCmd.AddCommand(goBasicCmd)
}

type Kubernetes struct {
	name string
	version string
	users []string
	nodeNumber int
}

func (k *Kubernetes) PrintUsers() string {
	return fmt.Sprintf("Users: %v", k.users)
}

func (k *Kubernetes) AddUser(user string) {
	k.users = append(k.users, user)
}

func (k *Kubernetes) RemoveUser(user string) {
	for i, u := range k.users {
		if u == user {
			k.users = append(k.users[:i], k.users[i+1:]...)
		}
	}
}