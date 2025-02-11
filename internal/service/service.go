package service

import (
	"context"
	"errors"
	"fmt"
	"promoter/internal/consts"
	factories "promoter/internal/factories/registry"
	"promoter/internal/types"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
)

type Service struct {
    project string
    name string
    env string
    config map[string]interface{}
    projectConfig *utils.Config 
}

func NewService(project, name, env string, config *utils.Config) (*Service, error) {

	serviceConfig, err := utils.FindService(config, name)
	if err != nil {
		return nil,err
	}

    return &Service{
        project,name,env,serviceConfig,config,
    }, nil
}

func (s *Service) Process(ctx context.Context, tag, region string, changeLog *[]types.ServiceChanges, interactive bool) error {
	repoName, err := utils.GetImageRepository(s.project, s.name, s.env)
	if err != nil {
		return err
	}

	registryFactory := &factories.RegistryFactory{}
	newTag := ""

    // TODO: take type from image names after implementing more registries
	registryClient, err := registryFactory.InitializeRegistry(ctx, consts.ECR, region)
	if err != nil {
		return err
	}

	if tag != consts.EmptyString {
		if err := registryClient.ImageExists(ctx, repoName, tag); err != nil {
			return err
		}
		newTag = tag
	} else {

		latestImage, err := registryClient.GetLatestImage(ctx, repoName)
		if err != nil {
			return err
		}
		newTag = latestImage.ImageTags[len(latestImage.ImageTags)-1]
	}

	didChange, err := s.changeTag(newTag, interactive)
	if err != nil {
		return err
	}

	if didChange {
		*changeLog = append(*changeLog, types.ServiceChanges{
			Name:   s.name,
			NewTag: newTag,
		})
	}
	return nil
}

func (s *Service) changeTag(tag string, interactive bool) (bool, error) {
	imageTagKey := utils.GetImageTagKey()

	if _, ok := s.config[imageTagKey]; !ok {
		return false, errors.New("image tag not found in the service's fields")
	}

	if s.config[imageTagKey] == tag {
		fmt.Printf("Service %s is already at latest / input tag \n", s.name)
		return false, nil
	}

	var change bool = true

	if interactive {

		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("Update %s to %s?", s.name, tag),
			Default: false,
		}
		if err := survey.AskOne(confirmPrompt, &change); err != nil {
			return false, err
		}
	} else {
		fmt.Printf("Updating service %s to %s\n", s.name, tag)
	}

	if !change {
		return false, nil
	}

	s.config[imageTagKey] = tag
    if err := utils.WriteToProjectFile(s.project, s.env, s.projectConfig); err != nil {
		return false, err
	}

	return true, nil
}
