package data

import (
	"errors"
)

type Config map[interface{}]interface{}

func GetApplications(project string, env string, projectFilePath string, manifestRepoRoot string) ([]string, error) {

	config, _, err := GetProjectConfig(project, env, projectFilePath, manifestRepoRoot)
	if err != nil {
		return nil, err
	}

	if _, ok := (*config)["applications"]; !ok {
		return nil, errors.New("application field not found in values file")
	}

	applications, ok := (*config)["applications"].([]interface{})
	if !ok {
		return nil, errors.New("applications field is not a list")
	}

	var services []string

	for _, app := range applications {
		appMap, ok := app.(map[string]interface{})

		name, ok := appMap["name"].(string)
		if !ok {
			continue
		}

		services = append(services, name)
	}
	return services, nil
}
