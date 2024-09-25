package commands

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetServicesCmd(cmd *cobra.Command) {
	env, _ := cmd.Flags().GetString("env")
	project, _ := cmd.Flags().GetString("project")
	projectFile, _ := cmd.Root().PersistentFlags().GetString("project-file")

	if projectFile == "" && (env == "" || project == "") {
		fmt.Print("You Need to either provide both env and project flags, or the project-file path")
		return
	}

	services, err := utils.GetServicesNames(project, env, projectFile)
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, app := range services {
		fmt.Println("* " + app)
	}
}
