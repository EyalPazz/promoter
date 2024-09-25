package utils

import (
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
