package commands

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetServicesCmd(cmd *cobra.Command) {

	env, _ := cmd.Flags().GetString("env")
	project, _ := cmd.Flags().GetString("project")

    project, _ , err := utils.ValidateProjectAttributes(project, "")

    if project == "" {
		fmt.Print(err)
		return
    }


	services, err := utils.GetServicesNames(project, env)
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, app := range services {
		fmt.Println("* " + app)
	}
}
