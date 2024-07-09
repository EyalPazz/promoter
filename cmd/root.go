package cmd

import (
	"context"
	"fmt"
	"os"
	"promoter/helpers/data"
	"promoter/helpers/manipulations"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	project     string
	service     string
	imageRepo   string
	env         string
	tag         string
	projectFile string

	rootCmd = &cobra.Command{
		Use:   "promoter",
		Short: "promoter is a CLI tool to easily deploy services",
		Long:  `promoter is a CLI tool to easily deploy services across different environments`,
		Run: func(cmd *cobra.Command, args []string) {

			passphrase, err := cmd.Flags().GetBool("passphrase")
			if err != nil {
				fmt.Print(err)
				return
			}
			data.RefreshRepo(passphrase)

			if region == "" {
				region = viper.GetString("region")
			}

			if region == "" {
				fmt.Println("Error: region must be specified either as flags or in the config file")
				return
			}

			ctx := context.Background()
			// IMPORTANT: Notice the convention for registry names
			var repoName string
			if imageRepo == "" {
				repoName = data.GetImageRepository(project, service)
			} else {
				repoName = imageRepo
			}
			if tag == "" {
				latestImage, err := data.GetLatestImage(ctx, repoName, region)
				if err != nil {
					fmt.Print(err)
					return
				}
				tag = latestImage.ImageTags[0]
			} else if err := data.ImageExists(ctx, repoName, tag, region); err != nil {
				fmt.Println(err)
				return
			}

			err = manipulations.ChangeServiceTag(project, service, env, tag, projectFile, viper.GetString("manifestRepoRoot"))
			if err != nil {
				fmt.Print(err)
				return
			}

			// passphraseFlag, err := cmd.PersistentFlags().GetBool("passphrase")
			// if err != nil {
			// 	fmt.Print(err)
			// }
			//
			// err = manipulations.CommitRepoChange(project, service, env, tag, passphraseFlag)
			// if err != nil {
			// 	fmt.Print(err)
			// }

		},
	}
)

// Execute executes the root command.
func Execute() error {
	Extend(rootCmd)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.promoter.yaml)")
	rootCmd.Flags().StringVar(&project, "project", "", "Project name (required)")
	rootCmd.Flags().StringVar(&service, "service", "", "Service name (required)")
	rootCmd.Flags().StringVar(&env, "env", "", "Environment name (required)")
	rootCmd.Flags().StringVar(&imageRepo, "image-repository", "", "Image repository name")
	rootCmd.Flags().StringVar(&tag, "tag", "", "Tag name")
	rootCmd.Flags().StringVar(&projectFile, "project-file", "", "Project File")
	rootCmd.PersistentFlags().Bool("passphrase", false, "Whether or not to prompt for ssh key passphrase")

	rootCmd.MarkFlagRequired("project")
	rootCmd.MarkFlagRequired("service")
	rootCmd.MarkFlagRequired("env")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".promoter")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error Reading Config File")
	}
}
