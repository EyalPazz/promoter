package cmd

import (
	"fmt"
	"promoter/helpers/data"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	region string
)

func Extend(rootCmd *cobra.Command) {

	rootCmd.AddCommand(RefreshManifestRepoCmd)

	GetApplicationsCmd.Flags().String("project", "", "Project name (required)")
	GetApplicationsCmd.Flags().String("service", "", "Service name (required)")
	GetApplicationsCmd.Flags().String("env", "", "Environment name (required)")
	rootCmd.AddCommand(GetApplicationsCmd)
}

var RefreshManifestRepoCmd = &cobra.Command{
	Use:   "refresh-manifest",
	Short: "Refresh The Manifest Repo",
	Long:  "Refresh The Manifest Repo",
	Run: func(cmd *cobra.Command, args []string) {
		passphraseFlag, err := cmd.Root().PersistentFlags().GetBool("passphrase")
		if err != nil {
			fmt.Print(err)
			return
		}
		data.RefreshRepo(passphraseFlag)
	},
}

var GetApplicationsCmd = &cobra.Command{
	Use:   "get-applications",
	Short: "Get All Applications in a certain project",
	Long:  "Get All Applications in a certain project",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		project, _ := cmd.Flags().GetString("project")
		projectFile, _ := cmd.Root().PersistentFlags().GetString("project-file")

		if projectFile == "" && (env == "" || project == "") {
			fmt.Print("You Need to either provide both env and project flags, or the project-file path")
			return
		}

		applications, err := data.GetApplicationsNames(project, env, projectFile, viper.GetString("manifestRepoRoot"))
		if err != nil {
			fmt.Print(err)
			return
		}

		for _, app := range applications {
			fmt.Println("* " + app)
		}
	},
}
