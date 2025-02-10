package utils

import (
	"fmt"
	"os"
	"promoter/internal/types"
	"slices"

	"github.com/spf13/viper"
)

const imageTagKey = "imageTag"

func GetImageTagKey() string {

	configImageTag := viper.GetString("ImageTag")
	if configImageTag != "" {
		return configImageTag
	}
	return imageTagKey
}

func ValidateProjectAttributes(project, region string) (string, string, error) {
	project = firstNonEmpty(project, viper.GetString("project-name"))
	region = firstNonEmpty(region, viper.GetString("region"))

	// Validate required fields
	if project == "" {
		return "", "", fmt.Errorf("project name must be specified either as a flag or in the config file")
	}
	if region == "" {
		return "", "", fmt.Errorf("region must be specified either as a flag or in the config file")
	}

	return project, region, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func GetConfig() (*types.Config, error) {

	config, ok := viper.Get("config").(types.Config)

	if !ok {
		return nil, fmt.Errorf("error: config structure is invalid")
	}

	return &config, nil
}

func GetGitProviderToken() string {
    token := os.Getenv("GIT_PROVIDER_TOKEN")
    if token == "" {
        fmt.Println("warning: GIT_PROVIDER_TOKEN is undefined")
    }
    return token
}

func ShouldCreatePR(env string) bool {
    config, _ := GetConfig()
     
    return config.PullRequests.Enabled && ((config.PullRequests.Envs == nil) || slices.Contains(config.PullRequests.Envs, env))
}
