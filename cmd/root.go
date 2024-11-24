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
	config      Config
	profile     string

	rootCmd = &cobra.Command{
		Use:   "promoter",
		Short: "promoter is a CLI tool to easily deploy services",
		Long:  `promoter is a CLI tool to easily deploy services across different environments`,
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

	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "Configuration profile to use")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.promoter.yaml)")
	rootCmd.PersistentFlags().Bool("passphrase", false, "Whether or not to prompt for ssh key passphrase")

	rootCmd.Flags().StringVar(&services, "services", "", "Services  (separeted by a comma)")
	rootCmd.Flags().StringVar(&project, "project", "", "Project name (required)")
	rootCmd.Flags().StringVar(&env, "env", "", "Environment name (required)")
	rootCmd.Flags().StringVar(&tag, "tag", "", "Specific image tag to promote (or revert) to (Only Supported With One Service)")

	rootCmd.MarkFlagRequired("env")
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

	viper.Set("project-name", selectedProfile.ProjectName)
	viper.Set("region", selectedProfile.Region)

    fmt.Println("All loaded settings:")
    for key, value := range viper.AllSettings() {
        fmt.Printf("%s: %v\n", key, value)
    }
}
