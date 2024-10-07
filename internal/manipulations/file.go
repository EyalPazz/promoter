package manipulations

import (
	"errors"
	"fmt"
	"promoter/internal/utils"
)

func ChangeServiceTag(project string, service string, env string, tag string, projectFilePath string) (bool, error) {

	config, err := utils.GetProjectConfig(project, env, projectFilePath)
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
		return false, errors.New("Image Tag Not found in the service's fields")
	}

	if err = utils.WriteToProjectFile(project, env, projectFilePath, config); err != nil {
		return false, err
	}

	return true, nil
}
