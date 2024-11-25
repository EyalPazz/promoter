package manipulations

import (
	"errors"
	"fmt"
	"promoter/internal/utils"
)

func ChangeServiceTag(project, service, env, tag string) (bool, error) {

	config, err := utils.GetProjectConfig(project, env)
	if err != nil {
		return false, err
	}

	app, err := utils.FindService(config, service)
	if err != nil {
		return false, err
	}

	imageTagKey := utils.GetImageTagKey()

	if app[imageTagKey] == tag {
		fmt.Printf("Service %s is already at latest tag \n", service)
		return false, nil
	}

	if _, ok := app[imageTagKey]; ok {
		app[imageTagKey] = tag
	} else {
		return false, errors.New("image tag not found in the service's fields")
	}

	if err = utils.WriteToProjectFile(project, env, config); err != nil {
		return false, err
	}

	fmt.Printf("Updating service %s to tag %s\n", service, tag)

	return true, nil
}
