package cmd

import (
	"promoter/cmd/commands"
	"promoter/internal/consts"

	"github.com/spf13/cobra"
)

func Extend(rootCmd *cobra.Command) error {

	rootCmd.AddCommand(RefreshManifestRepoCmd)

	rootCmd.AddCommand(GetServicesCmd)

	GetServicesCmd.Flags().String(consts.Env, consts.EmptyString, consts.EnvFDesc)

	if err := GetServicesCmd.MarkFlagRequired(consts.Env); err != nil {
		return err
	}

	rootCmd.AddCommand(RevertProjectCmd)

	RevertProjectCmd.Flags().String(consts.Env, consts.EmptyString, consts.EnvFDesc)
	RevertProjectCmd.Flags().Int(consts.Since, 7, consts.SinceFDesc)

	if err := RevertProjectCmd.MarkFlagRequired(consts.Env); err != nil {
		return err
	}

	rootCmd.AddCommand(GetProfileCmd)
	GetProfileCmd.Flags().Bool(consts.All, false, consts.AllFDesc)

	GetProfileCmd.AddCommand(AddProfileCmd)

	return nil

}

var RefreshManifestRepoCmd = &cobra.Command{
	Use:   consts.RefreshManifestCmd,
	Short: consts.RefreshManifestShort,
	Long:  consts.RefreshManifestLong,
	Run: func(cmd *cobra.Command, args []string) {
		commands.RefreshManifestRepoCmd(cmd)
	},
}

var RevertProjectCmd = &cobra.Command{
	Use:   consts.RevertProjectCmd,
	Short: consts.RevertProjectShort,
	Long:  consts.RevertProjectLong,
	Run: func(cmd *cobra.Command, args []string) {
		commands.RevertProject(cmd)
	},
}

var GetServicesCmd = &cobra.Command{
	Use:   consts.GetServicesCmd,
	Short: consts.GetServicesShort,
	Long:  consts.GetServicesLong,
	Run: func(cmd *cobra.Command, args []string) {
		commands.GetServicesCmd(cmd)
	},
}

var GetProfileCmd = &cobra.Command{
	Use:   consts.GetProfileCmd,
	Short: consts.GetProfileShort,
	Long:  consts.GetProfileLong,
	Run: func(cmd *cobra.Command, args []string) {
		commands.GetProfile(cmd)
	},
}

var AddProfileCmd = &cobra.Command{
	Use:   consts.AddProfileCmd,
	Short: consts.AddProfileShort,
	Long:  consts.AddProfileLong,
	Run: func(cmd *cobra.Command, args []string) {
		commands.AddProfile(cmd)
	},
}
