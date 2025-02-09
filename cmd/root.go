package cmd

import (
	"fmt"
	"os"

	"promoter/cmd/commands"
	"promoter/internal/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// Flags
	region  string
	profile string

	// Misc
	config      types.Config
	showVersion bool
	Version     string = "dev"

	rootCmd = &cobra.Command{
		Use:   "promoter",
		Short: "promoter is a CLI tool to easily deploy services",
		Long:  `promoter is a CLI tool to easily deploy services across different environments`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("%s\n", Version)
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			project, _ := cmd.Flags().GetString("project")
			services, _ := cmd.Flags().GetString("services")
			env, _ := cmd.Flags().GetString("env")
			tag, _ := cmd.Flags().GetString("tag")
			region, _ := cmd.Flags().GetString("region")
			commands.RootCmd(cmd, region, services, tag, project, env)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	if err := Extend(rootCmd); err != nil {
		return err
	}

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent Flags
	rootCmd.PersistentFlags().Bool("passphrase", false, "Whether or not to prompt for ssh key passphrase")

	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "Ask for confirmation in each change")

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version number")

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Configuration profile to use")

	rootCmd.PersistentFlags().String("services", "", "Services  (separeted by a comma)")
	rootCmd.PersistentFlags().String("project", "", "Project name (required)")

	// Root CMD Glags
	rootCmd.Flags().String("tag", "", "Specific image tag to promote (or revert) to (Only Supported With One Service)")
	rootCmd.Flags().String("env", "", "Environment name (required)")

	rootCmd.Flags().StringVar(&region, "region", "", "AWS Region for repository")

	err := rootCmd.MarkFlagRequired("env")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".promoter")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Print(fmt.Errorf("fatal error reading config file: %w", err))
		os.Exit(1)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Print(fmt.Errorf("unable to decode into config struct: %w", err))
		os.Exit(1)
	}

	selectedProfile, exists := config.Profiles[profile]
	if !exists {
		fmt.Print(fmt.Errorf("profile '%s' not found in configuration", profile))
		os.Exit(1)
	}

	viper.Set("config", config)

	if region == "" {
		region = selectedProfile.Region
	}

	viper.Set("project-name", selectedProfile.ProjectName)
	viper.Set("region", region)

}
