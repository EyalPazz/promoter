package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func GetProjectFile(repoPath string, project string, env string, manifestRepoRoot string) (string, error) {
	fileExtensions := []string{".yaml", ".yml"}
	for _, ext := range fileExtensions {
		projectFile := filepath.Join(repoPath, manifestRepoRoot, project, env, "values"+ext)
		if FileExists(projectFile) {
			return projectFile, nil
		}
	}
	return "", fmt.Errorf("Project File Does Not exist")
}

func GetProjectConfig(project string, env string, projectFilePath string, manifestRepoRoot string) (*Config, string, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return nil, "", fmt.Errorf("Error getting repository path: %s\n", err)
	}

	if projectFilePath == "" {
		projectFilePath, err = GetProjectFile(repoPath, project, env, manifestRepoRoot)
		if err != nil {
			return nil, "", err
		}
	}

	yamlFile, err := os.ReadFile(projectFilePath)
	if err != nil {
		return nil, "", err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, "", fmt.Errorf("Error unmarshalling YAML: %v", err)

	}

	return &config, projectFilePath, nil

}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func FindApplication(config *Config, service string) (map[string]interface{}, error) {
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
