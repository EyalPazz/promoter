package manipulations

import (
	"errors"
	"fmt"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
)

func ChangeServiceTag(project, service, env, tag string, interactive bool) (bool, error) {

	config, err := utils.GetProjectConfig(project, env)
	if err != nil {
		return false, err
	}

	app, err := utils.FindService(config, service)
	if err != nil {
		return false, err
	}

	imageTagKey := utils.GetImageTagKey()

	if _, ok := app[imageTagKey]; ok {
		app[imageTagKey] = tag
	} else {
		return false, errors.New("image tag not found in the service's fields")
	}

	if app[imageTagKey] == tag {
		fmt.Printf("Service %s is already at latest tag \n", service)
		return false, nil
	}

	var change bool = true

	if interactive {

		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("Update %s to %s?", service, tag),
			Default: false,
		}
		if err := survey.AskOne(confirmPrompt, &change); err != nil {
			return false, err
		}
	} else {
		fmt.Printf("Updating service %s to %s\n", service, tag)
	}

	if !change {
		return false, nil
	}

	if err = utils.WriteToProjectFile(project, env, config); err != nil {
		return false, err
	}

	return true, nil
}
