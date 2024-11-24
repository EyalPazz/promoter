package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func GetProjectConfig(project, env string) (*Config, error) {

	var err error

    projectFilePath, err := GetProjectFile(project, env, false)
    if err != nil {
        return nil, err
    }

	yamlFile, err := os.ReadFile(projectFilePath)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)

	}

	return &config, nil

}

func GetProjectFile(project string, env string, repoScoped bool) (string, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return "", fmt.Errorf("error getting repository path: %s", err)
	}

	fileExtensions := []string{".yaml", ".yml"}

	for _, ext := range fileExtensions {
		projectFile := filepath.Join(repoPath, viper.GetString("manifest-repo-root"), project, env, "values"+ext)
		if FileExists(projectFile) {
			if repoScoped {
				return filepath.Join(viper.GetString("manifest-repo-root"), project, env, "values"+ext), nil
			}
			return projectFile, nil
		}
	}

	return "", fmt.Errorf("project File Does Not exist")
}

func WriteToProjectFile(project, env string, config *Config) error {

    projectFilePath, err := GetProjectFile(project, env, false)
		if err != nil {
			return err
		}

	updatedYAML, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err = os.WriteFile(projectFilePath, updatedYAML, 0644); err != nil {
		return err
	}

	return nil
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
