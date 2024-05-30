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

// type Config struct {
//     Properties map[string]ServiceConfig
// }

const imageTagKey = "imageTag"

func ChangeServiceTag(project string, service string, env string, tag string) error {

	repoPath, err := data.GetRepoPath()
	if err != nil {
		fmt.Printf("Error getting repository path: %s\n", err)
		return err
	}

	// TODO: ADD .yaml / .yml support
	projectFile := filepath.Join(repoPath, project+"/"+service+"/"+"values-"+env+".yaml")
	yamlFile, err := os.ReadFile(projectFile)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return err
	}

	var config ServiceConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return err
	}

	if _, ok := config[imageTagKey]; ok {
		config[imageTagKey] = tag
	} else {
		return errors.New("Image Tag Not found in values file")
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