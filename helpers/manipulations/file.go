package manipulations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"promoter/helpers/data"

	"gopkg.in/yaml.v3"
)

type ServiceConfig map[string]interface{}

type Config struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

func ChangeServiceTag(fileName string, service string, env string, tag string ) error {

    repoPath, err := data.GetRepoPath()
	if err != nil {
		fmt.Printf("Error getting repository path: %s\n", err)
	    return err 
    }


    projectFile := filepath.Join(repoPath, fileName + ".yaml")
    yamlFile, err := os.ReadFile(projectFile)
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
    
    if serviceConfig, ok := config.Services[service]; ok {
		serviceConfig[env] = tag
		config.Services[service] = serviceConfig
	} else {
	fmt.Printf("Config: %+v\n", config)
		return errors.New("Service not found")
	}

    updatedYAML, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = os.WriteFile(projectFile, updatedYAML, 0644)
	if err != nil {
		return err
	}
    return nil
}
