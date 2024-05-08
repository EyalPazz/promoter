package cmd

import (
    "github.com/spf13/cobra"
    "promoter/helpers")

func Extend(rootCmd *cobra.Command) {
    rootCmd.AddCommand(helpers.ManifestRepoClone)
}

