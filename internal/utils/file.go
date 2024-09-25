package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)


func GetProjectConfig(project string, env string, projectFilePath string) (*Config, string, error) {

    var err error;

	if projectFilePath == "" {
        projectFilePath, err = GetProjectFile( project, env, false)
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

func FindService(config *Config, service string) (map[string]interface{}, error) {
	if _, ok := (*config)["applications"]; !ok {
		return nil, errors.New("application field not found in values file")
	}

	services, ok := (*config)["applications"].([]interface{})
	if !ok {
		return nil, errors.New("applications field is not a list")
	}

	for _, app := range services {
		appMap, ok := app.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := appMap["name"].(string)
		if !ok {
			continue
		}

		appType, ok := appMap["type"].(string)
		if !ok {
			continue
		}

		if name+"-"+appType == service {
			return appMap, nil
		}
	}
	return nil, fmt.Errorf("service with name '%s' not found", service)
}

func GetProjectFile(project string, env string, repoScoped bool) (string, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return "", fmt.Errorf("Error getting repository path: %s\n", err)
	}

	fileExtensions := []string{".yaml", ".yml"}

	for _, ext := range fileExtensions {
		projectFile := filepath.Join(repoPath, viper.GetString("manifestRepoRoot"), project, env, "values"+ext)
		if FileExists(projectFile) {
            if repoScoped {
                return filepath.Join(viper.GetString("manifestRepoRoot"), project, env, "values"+ext), nil
            }
			return projectFile, nil
		}
	}
	return "", fmt.Errorf("Project File Does Not exist")
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
