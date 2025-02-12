package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"promoter/internal/commands"
	"promoter/internal/consts"
	"promoter/internal/utils"
)

func RevertProject(cmd *cobra.Command) {
	env, _ := cmd.Flags().GetString(consts.Env)
	passphrase, _ := cmd.Flags().GetBool(consts.Passphrase)
	project, _ := cmd.Flags().GetString(consts.Project)
	since, _ := cmd.Flags().GetInt(consts.Since)

	project, _, err := utils.ValidateProjectAttributes(project, consts.EmptyString)

	if project == consts.EmptyString {
		fmt.Print(err)
		return
	}

	projectReverter := commands.NewProjectReverter(env, project)

	if err := projectReverter.Process(since, passphrase); err != nil {
		fmt.Println(err)
	}

	fmt.Println(consts.Success)
}
