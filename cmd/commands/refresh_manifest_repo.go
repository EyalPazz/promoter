package commands

import (
	"fmt"
	"promoter/internal/data"

	"github.com/spf13/cobra"
)

func RefreshManifestRepoCmd(cmd *cobra.Command) {
	passphraseFlag, err := cmd.Root().PersistentFlags().GetBool("passphrase")
	if err != nil {
		fmt.Print(err)
		return
	}

	if err = data.RefreshRepo(passphraseFlag); err != nil {
		fmt.Print(err)
	}
}
