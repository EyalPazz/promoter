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

func ValidateProjectAttributes(project string, region string)  (string, string, error){
	if project == "" {
		project = viper.GetString("project-name")
	}

	if project == "" {
		return "", "", fmt.Errorf("Error: project name must be specified either as flags or in the config file")
	}

	if region == "" {
		region = viper.GetString("region")
	}

	if region == "" {
		return "", "", fmt.Errorf("Error: region must be specified either as flags or in the config file")
	}

    return project, region, nil
}
