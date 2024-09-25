package manipulations

import (
	"errors"
	"fmt"
	"os"
	"promoter/internal/data"
	"promoter/internal/utils"

	"gopkg.in/yaml.v3"
)

func ChangeServiceTag(project string, service string, env string, tag string, projectFilePath string, manifestRepoRoot string) (error, bool) {

	config, projectFilePath, err := data.GetProjectConfig(project, env, projectFilePath, manifestRepoRoot)
	if err != nil {
		return err, false
	}

	app, err := data.FindService(config, service)
	if err != nil {
		return err, false
	}

	imageTagKey := utils.GetImageTagKey()

	if app[imageTagKey] == tag {
		fmt.Printf("Service %s is already at latest tag \n", service)
		return nil, false
	}

	if _, ok := app[imageTagKey]; ok {
		app[imageTagKey] = tag
	} else {
		return errors.New("Image Tag Not found in the service's fields"), false
	}

	updatedYAML, err := yaml.Marshal(config)
	if err != nil {
		return err, false
	}

	err = os.WriteFile(projectFilePath, updatedYAML, 0644)
	if err != nil {
		return err, false
	}
	return nil, true
}
