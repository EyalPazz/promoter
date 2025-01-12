package commands

import (
	"fmt"
	"promoter/internal/types"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddProfile(cmd *cobra.Command){
var profileName string
	if err := survey.AskOne(&survey.Input{
		Message: "Enter the profile name:",
	}, &profileName, survey.WithValidator(survey.Required)); err != nil {
        fmt.Printf("error: failed to get profile name: %v", err)
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
                Message: fmt.Sprintf("Profile '%s' already exists. Do you want to overwrite it?", profileName),
                Default: false,
            }, &overwrite); err != nil {
                fmt.Printf("Failed to confirm overwrite: %v", err)
            }
            if !overwrite {
                fmt.Println("Operation canceled.")
                return
            }
		}
	}

    var newProfile types.Profile;
    if err := survey.Ask([]*survey.Question{
		{
			Name:     "ProjectName",
			Prompt:   &survey.Input{Message: "Enter the project name:"},
			Validate: survey.Required,
		},
		{
			Name:     "Region",
			Prompt:   &survey.Input{Message: "Enter the region:"},
			Validate: survey.Required,
		},
	}, &newProfile); err != nil {
        fmt.Printf("error: Failed to get profile details: %v", err)
        return 
	}

    viper.Set("config",nil)
    viper.Set("project-name",nil)
    viper.Set("region",nil)

    viper.Set(profileName + ".project-name", newProfile.ProjectName)
    viper.Set(profileName + ".region", newProfile.Region)
    if err := viper.WriteConfig(); err != nil {
        fmt.Printf("Failed to write config: %v", err)
    }

	fmt.Printf("Profile '%s' saved successfully!\n", profileName)

}
