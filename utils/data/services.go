package data

import (
	"errors"
)

type Config map[interface{}]interface{}

func GetServiceImage(service string, project string, env string, projectFilePath string, manifestRepoRoot string) (string, error) {
	services, err := GetServices(project, env, projectFilePath, manifestRepoRoot)

	if err != nil {
		return "", err
	}

	for _, serviceConf := range services {
		serviceMap, ok := serviceConf.(map[string]interface{})
		name, ok := serviceMap["name"].(string)
		serviceType, ok := serviceMap["type"].(string)
		if !ok || name+"-"+serviceType != service {
			continue
		}

		image, ok := serviceMap["image"].(string)
		if !ok {
			return "", errors.New("Can't Find a valid image in servicelication's values")
		}
		return image, nil
	}

	return "", errors.New("Can't find the requested service in the project's values file")

}

func GetServices(project string, env string, projectFilePath string, manifestRepoRoot string) ([]interface{}, error) {

	config, _, err := GetProjectConfig(project, env, projectFilePath, manifestRepoRoot)
	if err != nil {
		return nil, err
	}

	if _, ok := (*config)["applications"]; !ok {
		return nil, errors.New("application field not found in values file")
	}

	services, ok := (*config)["applications"].([]interface{})
	if !ok {
		return nil, errors.New("applications field is not a list")
	}

	return services, nil
}

func GetServicesNames(project string, env string, projectFilePath string, manifestRepoRoot string) ([]string, error) {

	services, err := GetServices(project, env, projectFilePath, manifestRepoRoot)

	if err != nil {
		return nil, err
	}

	var serviceNames []string

	for _, app := range services {
		appMap, ok := app.(map[string]interface{})

		name, ok := appMap["name"].(string)
		if !ok {
			continue
		}

		appType, ok := appMap["type"].(string)
		if !ok {
			continue
		}

		serviceNames = append(serviceNames, name+"-"+appType)
	}
	return serviceNames, nil
}