package cmd

import (
	"promoter/cmd/commands"

	"github.com/spf13/cobra"
)

var (
	region string
)

func Extend(rootCmd *cobra.Command) {

	rootCmd.AddCommand(RefreshManifestRepoCmd)

	GetServicesCmd.Flags().String("project", "", "Project name (required)")
	GetServicesCmd.Flags().String("service", "", "Service name (required)")
	GetServicesCmd.Flags().String("env", "", "Environment name (required)")
	rootCmd.AddCommand(GetServicesCmd)

	rootCmd.AddCommand(RevertProjectCmd)
	RevertProjectCmd.Flags().String("project", "", "Project name (required)")
	RevertProjectCmd.Flags().String("service", "", "Service name (required)")
	RevertProjectCmd.Flags().String("env", "", "Environment name (required)")
	RevertProjectCmd.Flags().Int("since", 7, "Time interval to get revisions from (defaults to 7)")

}

var RefreshManifestRepoCmd = &cobra.Command{
	Use:   "refresh-manifest",
	Short: "Refresh The Manifest Repo",
	Long:  "Refresh The Manifest Repo",
	Run: func(cmd *cobra.Command, args []string) {
		commands.RefreshManifestRepoCmd(cmd)
	},
}

var RevertProjectCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert a service to a previous version",
	Long:  "Refresh The Manifest Repo",
	Run: func(cmd *cobra.Command, args []string) {
		commands.RevertProject(cmd)
	},
}

var GetServicesCmd = &cobra.Command{
	Use:   "get-services",
	Short: "Get All Services in a certain project",
	Long:  "Get All Services in a certain project",
	Run: func(cmd *cobra.Command, args []string) {
		commands.GetServicesCmd(cmd)
	},
}
