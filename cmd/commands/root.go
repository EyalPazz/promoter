package commands

import (
	"context"
	"fmt"
	factories "promoter/internal/factories/registry"
	"promoter/internal/manipulations"
	"promoter/internal/types"
	"promoter/internal/utils"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd(cmd *cobra.Command, region, services, tag, project, env, projectFile string) {
	passphrase, _ := cmd.Flags().GetBool("passphrase")

	var err error

	if region == "" {
		region = viper.GetString("region")
	}

	if region == "" {
		fmt.Println("Error: region must be specified either as flags or in the config file")
		return
	}

	if err = utils.RefreshRepo(passphrase); err != nil {
		fmt.Println(err)
		return
	}

	serviceList, err := getServices(services, project, env, projectFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if tag != "" && len(serviceList) > 1 {
		fmt.Println("Error: Image Tag Flag Only Supported With One Service")
		return
	}

	ctx := context.Background()

	var changeLog []types.ServiceChanges

	// TODO: Think about the trade-offs in making this async
	for _, service := range serviceList {
		if err := processService(ctx, project, service, env, tag, region, projectFile, &changeLog); err != nil {
			fmt.Println(err)
			if len(changeLog) > 0 {
				fmt.Println("Reverting Changes...")
				if err = manipulations.DiscardChanges(); err != nil {
					fmt.Println(err)
				}
			}
			return
		}

	}

	if len(changeLog) == 0 {
		fmt.Println("Nothing To Promote")
		return
	}

	commitMsg := utils.ComposeCommitMsg(&changeLog, env, project)

	if err := manipulations.HandleRepoActions(commitMsg, passphrase); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Success!")
}

func getServices(serviceStr, project, env, projectFile string) ([]string, error) {
	var serviceList []string
	var err error

	if serviceStr == "" {
		serviceList, err = utils.GetServicesNames(project, env, projectFile)
		if err != nil {
			return nil, fmt.Errorf("error retrieving service names: %v", err)
		}
	} else {
		serviceList = strings.Split(serviceStr, ",")
	}

	return serviceList, nil
}

func processService(ctx context.Context, project, service, env, tag, region, projectFile string, changeLog *[]types.ServiceChanges) error {
	repoName, err := utils.GetImageRepository(project, service, env, projectFile)
	if err != nil {
		return err
	}

	registryFactory := &factories.RegistryFactory{}
	newTag := ""

	ecrClient, err := registryFactory.InitializeRegistry(ctx, "ecr", region)
	if err != nil {
		return err
	}

	if tag != "" {
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

	didChange, err := manipulations.ChangeServiceTag(project, service, env, newTag, projectFile)
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
