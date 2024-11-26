package utils

import (
	"fmt"

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
