package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func GetProjectConfig(project, env, projectFilePath string) (*Config, error) {

	var err error

	if projectFilePath == "" {
		projectFilePath, err = GetProjectFile(project, env, false)
		if err != nil {
			return nil, err
		}
	}

	yamlFile, err := os.ReadFile(projectFilePath)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("Error unmarshalling YAML: %v", err)

	}

	return &config, nil

}

func GetProjectFile(project string, env string, repoScoped bool) (string, error) {
	repoPath, err := GetRepoPath()
	if err != nil {
		return "", fmt.Errorf("Error getting repository path: %s\n", err)
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

	return "", fmt.Errorf("Project File Does Not exist")
}

func WriteToProjectFile(project, env, projectFilePath string, config *Config) error {
	var err error

	if projectFilePath == "" {
		projectFilePath, err = GetProjectFile(project, env, false)
		if err != nil {
			return err
		}
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
