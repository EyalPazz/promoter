package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func RefreshManifestRepoCmd(cmd *cobra.Command) {
	passphraseFlag, err := cmd.Root().PersistentFlags().GetBool(consts.Passphrase)
	if err != nil {
		fmt.Print(err)
		return
	}

	if err = utils.RefreshRepo(passphraseFlag); err != nil {
		fmt.Print(err)
	}
}
