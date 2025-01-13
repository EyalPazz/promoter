package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/types"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddProfile(cmd *cobra.Command) {
	var profileName string
	if err := survey.AskOne(&survey.Input{
		Message: consts.EnterProfileName,
	}, &profileName, survey.WithValidator(survey.Required)); err != nil {
		fmt.Printf(consts.FailedToGetProfileName, err)
		return
	}

	config, err := utils.GetConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	for profile := range config.Profiles {
		if profile == profileName {
			var overwrite bool
			if err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf(consts.OverwriteProfile, profileName),
				Default: false,
			}, &overwrite); err != nil {
				fmt.Printf(consts.FailedToConfirmOverwrite, err)
			}
			if !overwrite {
				fmt.Println(consts.OperationCanceled)
				return
			}
		}
	}

	var newProfile types.Profile
	if err := survey.Ask([]*survey.Question{
		{
			Name:     consts.ProjectNameP,
			Prompt:   &survey.Input{Message: consts.EnterProjectName},
			Validate: survey.Required,
		},
		{
			Name:     consts.RegionP,
			Prompt:   &survey.Input{Message: consts.EnterRegion},
			Validate: survey.Required,
		},
	}, &newProfile); err != nil {
		fmt.Printf(consts.FailedToGetProfileDetails, err)
		return
	}

	viper.Set(consts.Config, nil)
	viper.Set(consts.Region, nil)
	viper.Set(consts.ProjectName, nil)

	viper.Set(profileName+consts.Dot+consts.ProjectName, newProfile.ProjectName)
	viper.Set(profileName+consts.Dot+consts.Region, newProfile.Region)
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf(consts.FailedToWriteConfig, err)
	}

	fmt.Printf("Profile '%s' saved successfully!\n", profileName)

}
