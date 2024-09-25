package commands

import (
	"fmt"
	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func RefreshManifestRepoCmd(cmd *cobra.Command) {
	passphraseFlag, err := cmd.Root().PersistentFlags().GetBool("passphrase")
	if err != nil {
		fmt.Print(err)
		return
	}

	if err = utils.RefreshRepo(passphraseFlag); err != nil {
		fmt.Print(err)
	}
}
