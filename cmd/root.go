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
			if tag == "" {
				// Notice That the convention for registry
				latestImage, err := data.GetLatestImage(ctx, repositoryName, region)
				if err != nil {
					fmt.Print(err)
					return
				}
				tag = latestImage.ImageTags[0]
			} else if err := data.ImageExists(ctx, repositoryName, tag, region); err != nil {
				fmt.Println(err)
				return
			}

			err = manipulations.ChangeServiceTag(project, service, env, tag, projectFile)
			if err != nil {
				fmt.Print(err)
				return
			}

			passphraseFlag, err := cmd.PersistentFlags().GetBool("passphrase")
			if err != nil {
				fmt.Print(err)
			}

			err = manipulations.CommitRepoChange(project, service, env, tag, passphraseFlag)
			if err != nil {
				fmt.Print(err)
			}

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
