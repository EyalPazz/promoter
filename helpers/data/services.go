package data

import (
	"errors"
)

type Config map[interface{}]interface{}

func GetServiceImage(service string, project string, env string, projectFilePath string, manifestRepoRoot string) (string, error) {
	applications, err := GetApplications(project, env, projectFilePath, manifestRepoRoot)

	if err != nil {
		return "", err
	}

	for _, app := range applications {
		appMap, ok := app.(map[string]interface{})
		name, ok := appMap["name"].(string)
		appType, ok := appMap["type"].(string)
		if !ok || name+"-"+appType != service {
			continue
		}

		image, ok := appMap["image"].(string)
		if !ok {
			return "", errors.New("Can't Find a valid image in application's values")
		}
		return image, nil
	}

	return "", errors.New("Can't find the requested service in the project's values file")

}

func GetApplications(project string, env string, projectFilePath string, manifestRepoRoot string) ([]interface{}, error) {

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

	return applications, nil
}

func GetApplicationsNames(project string, env string, projectFilePath string, manifestRepoRoot string) ([]string, error) {

	applications, err := GetApplications(project, env, projectFilePath, manifestRepoRoot)

	if err != nil {
		return nil, err
	}

	var services []string

	for _, app := range applications {
		appMap, ok := app.(map[string]interface{})

		name, ok := appMap["name"].(string)
		if !ok {
			continue
		}

		appType, ok := appMap["type"].(string)
		if !ok {
			continue
		}

		services = append(services, name+"-"+appType)
	}
	return services, nil
}
