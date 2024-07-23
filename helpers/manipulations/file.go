package manipulations

import (
	"errors"
	"fmt"
	"os"
	"promoter/helpers/data"

	"gopkg.in/yaml.v3"
)

const imageTagKey = "imageTag"

func ChangeServiceTag(project string, service string, env string, tag string, projectFilePath string, manifestRepoRoot string) error {

	config, projectFilePath, err := data.GetProjectConfig(project, env, projectFilePath, manifestRepoRoot)
	if err != nil {
		return err
	}

	app, err := data.FindApplication(config, service)
	if err != nil {
		return err
	}

	if app[imageTagKey] == tag {
		fmt.Println("Service is already at latest/input tag")
		return nil
	}

	if _, ok := app[imageTagKey]; ok {
		app[imageTagKey] = tag
	} else {
		return errors.New("Image Tag Not found in the service's fields")
	}

	updatedYAML, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(projectFilePath, updatedYAML, 0644)
	if err != nil {
		return err
	}
	return nil
}
