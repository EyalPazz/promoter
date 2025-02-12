package commands

import (
	"fmt"
	"promoter/internal/commands"
	"promoter/internal/consts"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func RootCmd(cmd *cobra.Command, region, services, tag, project, env string) {
	passphrase, _ := cmd.Flags().GetBool(consts.Passphrase)
	interactive, _ := cmd.Flags().GetBool(consts.Interactive)

	var err error

	project, region, err = utils.ValidateProjectAttributes(project, region)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = utils.RefreshRepo(passphrase); err != nil {
		fmt.Println(err)
		return
	}

	projectInstance, err := commands.NewProjectPromoter(services, env, project)

	if err != nil {
		fmt.Println(err)
		return
	}

	if tag != consts.EmptyString && len(*projectInstance.Services) > 1 {
		fmt.Println(consts.ImageTagFlagNotSupported)
		return
	}

	if err := projectInstance.Process(tag, region, interactive, passphrase); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(consts.Success)
}
