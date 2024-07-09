package manipulations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"promoter/helpers/data"

	"gopkg.in/yaml.v3"
)

type Config map[interface{}]interface{}

const imageTagKey = "imageTag"

func ChangeServiceTag(project string, service string, env string, tag string, projectFilePath string, manifestRepoRoot string) error {

	repoPath, err := data.GetRepoPath()
	if err != nil {
		fmt.Printf("Error getting repository path: %s\n", err)
		return err
	}

	if projectFilePath == "" {
		projectFilePath, err = getProjectFile(repoPath, project, env, manifestRepoRoot)
		if err != nil {
			return err
		}
	}

	yamlFile, err := os.ReadFile(projectFilePath)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return err
	}

	app, err := findApplication(&config, service)
	if err != nil {
		return err
	}

	if _, ok := app[imageTagKey]; ok {
		app[imageTagKey] = tag
	} else {
		return errors.New("Image Tag Not found in the service's fields")
	}

	updatedYAML, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = os.WriteFile(projectFilePath, updatedYAML, 0644)
	if err != nil {
		return err
	}
	return nil
}

func findApplication(config *Config, service string) (map[string]interface{}, error) {
	if _, ok := (*config)["applications"]; !ok {
		return nil, errors.New("application field not found in values file")
	}

	applications, ok := (*config)["applications"].([]interface{})
	if !ok {
		return nil, errors.New("applications field is not a list")
	}

	for _, app := range applications {
		appMap, ok := app.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := appMap["name"].(string)
		if !ok {
			continue
		}

		if name == service {
			return appMap, nil
		}
	}
	return nil, fmt.Errorf("service with name '%s' not found", service)
}

func getProjectFile(repoPath string, project string, env string, manifestRepoRoot string) (string, error) {
	fileExtensions := []string{".yaml", ".yml"}
	for _, ext := range fileExtensions {
		projectFile := filepath.Join(repoPath, manifestRepoRoot, project, env, "values"+ext)
		if FileExists(projectFile) {
			return projectFile, nil
		}
	}
	return "", fmt.Errorf("Project File Does Not exist")
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
