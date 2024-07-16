package cmd

import (
	"context"
	"fmt"
	"promoter/helpers/data"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	region string
)

func Extend(rootCmd *cobra.Command) {

	GetLatestImageCmd.Flags().StringVarP(&region, "region", "R", "", "AWS region")
	rootCmd.AddCommand(GetLatestImageCmd)

	rootCmd.AddCommand(RefreshManifestRepoCmd)

	GetApplicationsCmd.Flags().String("project", "", "Project name (required)")
	GetApplicationsCmd.Flags().String("service", "", "Service name (required)")
	GetApplicationsCmd.Flags().String("env", "", "Environment name (required)")
	rootCmd.AddCommand(GetApplicationsCmd)
}

var GetLatestImageCmd = &cobra.Command{
	Use:   "latest",
	Short: "Get Latest Image from desired repo",
	Long:  "Get Latest Image from desired repo",
	Run: func(cmd *cobra.Command, args []string) {

		if region == "" {
			region = viper.GetString("region")
		}

		if region == "" {
			fmt.Println("Error: region must be specified either as flags or in the config file")
			return
		}
		ctx := context.Background()
		imageDetail, err := data.GetLatestImage(ctx, data.GetImageRepository(project, service), region)
		if err != nil {
			fmt.Println("Error Getting Latest Image From Repo:", err)
			return
		}

		fmt.Println(imageDetail.ImageTags)
	},
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
    Use: "get-applications",
    Short: "Get All Applications in a certain project", 
	Long:  "Get All Applications in a certain project",
	Run: func(cmd *cobra.Command, args []string) {
		env, err := cmd.Flags().GetString("env")
		project, err := cmd.Flags().GetString("project")
		projectFile, err := cmd.Root().PersistentFlags().GetString("project-file")
		if err != nil {
			fmt.Print(err)
            return
		}

        applications, err := data.GetApplications(project, env, projectFile, viper.GetString("manifestRepoRoot"))
		if err != nil {
			fmt.Print(err)
            return
		}

        for _, app := range applications {
            fmt.Println("* " + app)
        }
	},

}
