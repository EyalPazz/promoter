package commands

import (
	"fmt"
	"promoter/internal/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetProfile(cmd *cobra.Command) {

	profile, _ := cmd.Flags().GetString("profile")
	all, _ := cmd.Flags().GetBool("all")

	config, ok := viper.Get("config").(types.Config)

	if !ok {
		fmt.Println("error: config structure is invalid")
		return
	}

	fmt.Printf("Active Profile is: %s \n", profile)
	fmt.Printf("Project Name: %s \n", config.Profiles[profile].ProjectName)
	fmt.Printf("Region: %s \n", config.Profiles[profile].Region)

	if all {
		fmt.Println("\nOther Profile Names:")
		for name := range config.Profiles {
			fmt.Printf("* %s \n", name)
		}
	}
}
