package cmd

import (
	"context"
	"fmt"
	"promoter/helpers/data"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repositoryName string
	region         string
)

func Extend(rootCmd *cobra.Command) {

	GetLatestImageCmd.Flags().StringVarP(&repositoryName, "repository", "r", "", "ECR repository name")
	GetLatestImageCmd.Flags().StringVarP(&region, "region", "R", "", "AWS region")
	rootCmd.AddCommand(GetLatestImageCmd)

	rootCmd.AddCommand(RefreshManifestRepoCmd)
}

var GetLatestImageCmd = &cobra.Command{
	Use:   "latest",
	Short: "Get Latest Image from desired repo",
	Long:  "Get Latest Image from desired repo",
	Run: func(cmd *cobra.Command, args []string) {

		if repositoryName == "" {
			repositoryName = viper.GetString("repository")
		}
		if region == "" {
			region = viper.GetString("region")
		}

		if repositoryName == "" || region == "" {
			fmt.Println("Error: repository and region must be specified either as flags or in the config file")
			return
		}
		ctx := context.Background()
		imageDetail, err := data.GetLatestImage(ctx, repositoryName, region)
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
		}
		data.RefreshRepo(passphraseFlag)
	},
}
