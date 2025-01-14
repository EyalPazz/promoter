package commands

import (
	"context"
	"fmt"
	"promoter/internal/consts"
	factories "promoter/internal/factories/registry"
	"promoter/internal/manipulations"
	"promoter/internal/types"
	"promoter/internal/utils"
	"strings"

	"github.com/spf13/cobra"
)

func RootCmd(cmd *cobra.Command, region, services, tag, project, env string) {
	passphrase, _ := cmd.Flags().GetBool(consts.Passphrase)
	interactive, _ := cmd.Flags().GetBool(consts.Interactive)

	var err error

	project, region, err = utils.ValidateProjectAttributes(project, region)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = utils.RefreshRepo(passphrase); err != nil {
		fmt.Println(err)
		return
	}

	serviceList, err := getServices(services, project, env)
	if err != nil {
		fmt.Println(err)
		return
	}

	if tag != consts.EmptyString && len(serviceList) > 1 {
		fmt.Println(consts.ImageTagFlagNotSupported)
		return
	}

	ctx := context.Background()

	var changeLog []types.ServiceChanges

	// TODO: Think about the trade-offs in making this async
	for _, service := range serviceList {
		if err := processService(ctx, project, service, env, tag, region, &changeLog, interactive); err != nil {
			fmt.Println(err)
			if len(changeLog) > 0 {
				fmt.Println(consts.RevertingChanges)
				if err = manipulations.DiscardChanges(); err != nil {
					fmt.Println(err)
				}
			}
			return
		}

	}

	if len(changeLog) == 0 {
		fmt.Println(consts.NothingToPromote)
		return
	}

	commitMsg := utils.ComposeCommitMsg(&changeLog, env, project)

	if err := manipulations.HandleRepoActions(commitMsg, passphrase); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(consts.Success)
}

func getServices(serviceStr, project, env string) ([]string, error) {
	var serviceList []string
	var err error

	if serviceStr == consts.EmptyString {
		serviceList, err = utils.GetServicesNames(project, env)
		if err != nil {
			return nil, fmt.Errorf(consts.ErrorRetrievingServiceNames, err)
		}
	} else {
		serviceList = strings.Split(serviceStr, consts.Comma)
	}

	return serviceList, nil
}

func processService(ctx context.Context, project, service, env, tag, region string, changeLog *[]types.ServiceChanges, interactive bool) error {
	repoName, err := utils.GetImageRepository(project, service, env)
	if err != nil {
		return err
	}

	registryFactory := &factories.RegistryFactory{}
	newTag := ""

	ecrClient, err := registryFactory.InitializeRegistry(ctx, consts.ECR, region)
	if err != nil {
		return err
	}

	if tag != consts.EmptyString {
		if err := ecrClient.ImageExists(ctx, repoName, tag); err != nil {
			return err
		}
		newTag = tag
	} else {
		// TODO: take type from image names after implementing more registries

		latestImage, err := ecrClient.GetLatestImage(ctx, repoName)
		if err != nil {
			return err
		}
		newTag = latestImage.ImageTags[len(latestImage.ImageTags)-1]
	}

	didChange, err := manipulations.ChangeServiceTag(project, service, env, newTag, interactive)
	if err != nil {
		return err
	}

	if didChange {

		*changeLog = append(*changeLog, types.ServiceChanges{
			Name:   service,
			NewTag: newTag,
		})
	}
	return nil
}
