package cmd

import (
	"fmt"
	"os"

	"promoter/cmd/commands"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	project     string
	services    string
	env         string
	projectFile string

	rootCmd = &cobra.Command{
		Use:   "promoter",
		Short: "promoter is a CLI tool to easily deploy services",
		Long:  `promoter is a CLI tool to easily deploy services across different environments`,
		Run: func(cmd *cobra.Command, args []string) {
			commands.RootCmd(cmd, region, services, project, env, projectFile)
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
	rootCmd.Flags().StringVar(&services, "services", "", "Services  (separeted by a comma)")
	rootCmd.Flags().StringVar(&env, "env", "", "Environment name (required)")
	rootCmd.PersistentFlags().StringVar(&projectFile, "project-file", "", "Project File")
	rootCmd.PersistentFlags().Bool("passphrase", false, "Whether or not to prompt for ssh key passphrase")

	rootCmd.MarkFlagRequired("project")
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
