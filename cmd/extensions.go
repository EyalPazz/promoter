package cmd

import (
	"fmt"
	"promoter/helpers/data"

	"github.com/spf13/cobra"
)

func Extend(rootCmd *cobra.Command) {
	rootCmd.AddCommand(ManifestRepoClone)
}

var ManifestRepoClone = &cobra.Command{
	Use:   "clone-manifest",
	Short: "Authenticate with Git provider using SSH key",
	Long:  "Authenticate with Git provider using SSH key to perform Git operations.",
	Run: func(cmd *cobra.Command, args []string) {

		cloneErr := data.CloneRepository()
		if cloneErr != nil {
			fmt.Println("Error cloning manifest repo:", cloneErr)
			return
		}

		fmt.Println("Successfully cloned manifest repo")
	},
}
