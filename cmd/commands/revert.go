package commands

import (
	// "fmt"
	// "promoter/internal/data"

	"promoter/internal/utils"

	"github.com/spf13/cobra"
)

func RevertService(cmd *cobra.Command) {
    utils.GetLatestRevisions("heroes-live-on", "staging", 7 ) 
}
