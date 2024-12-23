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
	tag         string
	region      string
	config      Config
	showVersion bool
	profile     string
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
			commands.RootCmd(cmd, region, services, tag, project, env)
		},
	}
)

type Profile struct {
	ProjectName string `mapstructure:"project-name"`
	Region      string `mapstructure:"region"`
}

type Config struct {
	GitName          string             `mapstructure:"git-name"`
	GitEmail         string             `mapstructure:"git-email"`
	SSHKey           string             `mapstructure:"ssh-key"`
	ManifestRepo     string             `mapstructure:"manifest-repo"`
	ManifestRepoRoot string             `mapstructure:"manifest-repo-root"`
	Profiles         map[string]Profile `mapstructure:",remain"`
}

// Execute executes the root command.
func Execute() error {
	Extend(rootCmd)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Show version number")

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Configuration profile to use")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.promoter.yaml)")
	rootCmd.PersistentFlags().Bool("passphrase", false, "Whether or not to prompt for ssh key passphrase")
	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "Ask for confirmation in each change")

	rootCmd.Flags().StringVar(&region, "region", "", "AWS Region for repository")
	rootCmd.Flags().StringVar(&services, "services", "", "Services  (separeted by a comma)")
	rootCmd.Flags().StringVar(&project, "project", "", "Project name (required)")
	rootCmd.Flags().StringVar(&env, "env", "", "Environment name (required)")
	rootCmd.Flags().StringVar(&tag, "tag", "", "Specific image tag to promote (or revert) to (Only Supported With One Service)")

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
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Print(fmt.Errorf("unable to decode into config struct: %w", err))
	}

	selectedProfile, exists := config.Profiles[profile]
	if !exists {
		fmt.Print(fmt.Errorf("profile '%s' not found in configuration", profile))
	}

	if region == "" {
		region = selectedProfile.Region
	}

	viper.Set("project-name", selectedProfile.ProjectName)
	viper.Set("region", region)

}
