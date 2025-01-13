package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func GetServicesCmd(cmd *cobra.Command) {

	env, _ := cmd.Flags().GetString(consts.Env)
	project, _ := cmd.Flags().GetString(consts.Project)
	passphrase, _ := cmd.Flags().GetBool(consts.Passphrase)

	project, _, err := utils.ValidateProjectAttributes(project, "")

	if err != nil {
		fmt.Print(err)
		return
	}

	if err = utils.RefreshRepo(passphrase); err != nil {
		fmt.Print(err)
		return
	}

	serviceAttributes, err := utils.GetServicesFields(project, env, consts.Name, consts.Type, consts.ImageTag)
	if err != nil {
		fmt.Print(err)
		return
	}

	for _, atts := range serviceAttributes {
		// TODO: Assert Types Before Print
		fmt.Printf("* %s-%s : %s \n", atts[consts.Name].(string), atts[consts.Type].(string), atts[consts.ImageTag].(string))
	}
}
