package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetProfile(cmd *cobra.Command) {

	profile, _ := cmd.Flags().GetString(consts.Profile)
	all, _ := cmd.Flags().GetBool(consts.All)

	config, err := utils.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(consts.ActiveProfileIs, profile)
	fmt.Printf(consts.ProjectNameP+consts.PlaceholderNewline, config.Profiles[profile].ProjectName)
	fmt.Printf(consts.RegionP+consts.PlaceholderNewline, config.Profiles[profile].Region)

	if all {
		fmt.Println(consts.OtherProfileNames)
		for name := range config.Profiles {
			fmt.Printf(consts.Asterisk+consts.PlaceholderNewline, name)
		}
	}
}
