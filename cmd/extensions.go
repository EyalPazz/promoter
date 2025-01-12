package cmd

import (
	"promoter/cmd/commands"

	"github.com/spf13/cobra"
)

func Extend(rootCmd *cobra.Command) error {

	rootCmd.AddCommand(RefreshManifestRepoCmd)

	rootCmd.AddCommand(GetServicesCmd)

	GetServicesCmd.Flags().String("project", "", "Project name")
	GetServicesCmd.Flags().String("service", "", "Service name (required)")
	GetServicesCmd.Flags().String("env", "", "Environment name (required)")
	if err := GetServicesCmd.MarkFlagRequired("env"); err != nil {
		return err
	}

	rootCmd.AddCommand(RevertProjectCmd)

	RevertProjectCmd.Flags().String("project", "", "Project name")
	RevertProjectCmd.Flags().String("service", "", "Service name (required)")
	RevertProjectCmd.Flags().String("env", "", "Environment name (required)")
	RevertProjectCmd.Flags().Int("since", 7, "Time interval to get revisions from (in days, defaults to 7)")

	if err := RevertProjectCmd.MarkFlagRequired("env"); err != nil {
		return err
	}

	rootCmd.AddCommand(GetProfileCmd)
	GetProfileCmd.Flags().Bool("all", false, "return all profile names")

	GetProfileCmd.AddCommand(AddProfileCmd)

	return nil

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

var GetProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Get Active (Or all services)",
	Long:  "Get Active (Or all services)",
	Run: func(cmd *cobra.Command, args []string) {
		commands.GetProfile(cmd)
	},
}

var AddProfileCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a config profile",
	Long:  "Add a config profile",
	Run: func(cmd *cobra.Command, args []string) {
		commands.AddProfile(cmd)
	},
}
