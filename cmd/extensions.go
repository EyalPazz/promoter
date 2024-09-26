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


	rootCmd.AddCommand(RevertServiceCmd)
	RevertServiceCmd.Flags().String("project", "", "Project name (required)")
	RevertServiceCmd.Flags().String("service", "", "Service name (required)")
	RevertServiceCmd.Flags().String("env", "", "Environment name (required)")
    
}

var RefreshManifestRepoCmd = &cobra.Command{
	Use:   "refresh-manifest",
	Short: "Refresh The Manifest Repo",
	Long:  "Refresh The Manifest Repo",
	Run: func(cmd *cobra.Command, args []string) {
		commands.RefreshManifestRepoCmd(cmd)
	},
}

var RevertServiceCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert a service to a previous version",
	Long:  "Refresh The Manifest Repo",
	Run: func(cmd *cobra.Command, args []string) {
		commands.RevertService(cmd)
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
