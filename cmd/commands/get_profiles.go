package commands

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetProfile(cmd *cobra.Command) {

	profile, _ := cmd.Flags().GetString("profile")
	all, _ := cmd.Flags().GetBool("all")

	config, err := utils.GetConfig()
	if err != nil {
		fmt.Println(err)
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
