package commands

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetServicesCmd(cmd *cobra.Command) {

	env, _ := cmd.Flags().GetString("env")
	project, _ := cmd.Flags().GetString("project")

	project, _, err := utils.ValidateProjectAttributes(project, "")

	if project == "" {
		fmt.Print(err)
		return
	}

	serviceAttributes, err := utils.GetServicesFields(project, env, "name", "type", "imageTag")
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, atts := range serviceAttributes {
		// TODO: Assert Types Before Print
		fmt.Printf("* %s-%s : %s \n", atts["name"].(string), atts["type"].(string), atts["imageTag"].(string))
	}
}
